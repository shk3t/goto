package query

import (
	"context"
	db "goto/src/database"
	"goto/src/model"
	"goto/src/utils"
	"time"

	"github.com/jackc/pgx/v5"
)

func readSolutionRow(row Scanable) *model.Solution {
	solution := model.Solution{}
	err := row.Scan(
		&solution.Id,
		&solution.UserId,
		&solution.TaskId,
		&solution.Status,
		&solution.Result,
		&solution.UpdatedAt,
	)
	if err != nil {
		return nil
	}
	return &solution
}

func readSolutionRows(rows pgx.Rows) map[int]model.Solution {
	solutionsByIds := map[int]model.Solution{}
	for rows.Next() {
		solution := readSolutionRow(rows)
		solutionsByIds[solution.Id] = *solution
	}
	return solutionsByIds
}

func readSolutionRowThenExtend(ctx context.Context, row pgx.Row) *model.Solution {
	solution := readSolutionRow(row)
	if solution == nil {
		return nil
	}

	files := getFilesBySolutions(ctx, []int{solution.Id})
	for _, sf := range files {
		solution.Files[sf.Name] = sf.Code
	}
	return solution
}

func readSolutionRowsThenExtend(ctx context.Context, rows pgx.Rows) []model.Solution {
	tasksByIds := readSolutionRows(rows)
	files := getFilesByTasks(ctx, utils.MapKeys(tasksByIds))
	for _, tf := range files {
		tasksByIds[tf.TaskId].Files[tf.Name] = tf.Path
	}
	return utils.MapValues(tasksByIds)
}

func GetSolution(ctx context.Context, id int) *model.Solution {
	row := db.ConnPool.QueryRow(ctx, "SELECT * FROM solution WHERE id = $1", id)
	return readSolutionRow(row)
}

func GetUserSolution(ctx context.Context, id int, userId int) *model.Solution {
	row := db.ConnPool.QueryRow(
		ctx,
		"SELECT * FROM solution WHERE id = $1 AND user_id = $2",
		id, userId,
	)
	return readSolutionRowThenExtend(ctx, row)
}

func GetUserSolutions(ctx context.Context, userId int) []model.Solution {
	rows, _ := db.ConnPool.Query(ctx, "SELECT * FROM solution WHERE user_id = $1", userId)
	return readSolutionRowsThenExtend(ctx, rows)
}

func getFilesBySolutions(ctx context.Context, solutionIds []int) []model.SolutionFile {
	files := []model.SolutionFile{}

	rows, _ := db.ConnPool.Query(
		ctx,
		"SELECT id, task_id, name, path FROM solution_file WHERE solution_id = ANY ($1)",
		solutionIds,
	)

	for rows.Next() {
		file := model.SolutionFile{}
		rows.Scan(&file.Id, &file.SolutionId, &file.Name, &file.Code)
		files = append(files, file)
	}

	return files
}

func saveSolutionFiles(ctx context.Context, tx pgx.Tx, s *model.Solution) {
	solutionFileEntries := [][]any{}
	for name, code := range s.Files {
		solutionFileEntries = append(solutionFileEntries, []any{s.Id, name, code})
	}
	tx.CopyFrom(
		ctx,
		pgx.Identifier{"solution_file"},
		[]string{"solution_id", "name", "code"},
		pgx.CopyFromRows(solutionFileEntries),
	)
}

func updateSolution(ctx context.Context, tx pgx.Tx, s *model.Solution) {
	tx.Exec(
		ctx, `
        UPDATE solution
        SET status = $1, result = $2, updated_at = $3
        WHERE id = $4`,
		s.Status, s.Result, s.UpdatedAt,
		s.Id,
	)
}

func SaveSolution(ctx context.Context, s *model.Solution) *model.Solution {
	tx, _ := db.ConnPool.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)
	s.UpdatedAt = time.Now()

	if s.Id == 0 {
		s.Status = "check"

		tx.QueryRow(
			ctx,
			"SELECT id from solution WHERE user_id = $1 AND task_id = $2",
			s.UserId, s.TaskId,
		).Scan(&s.Id)

		if s.Id == 0 {
			tx.QueryRow(
				ctx, `
                INSERT INTO solution (user_id, task_id, updated_at)
                VALUES ($1, $2, $3)
                RETURNING id`,
				s.UserId, s.TaskId, s.UpdatedAt,
			).Scan(&s.Id)
		} else {
			updateSolution(ctx, tx, s)
		}
	} else {
		updateSolution(ctx, tx, s)
	}

	tx.Exec(ctx, "DELETE FROM solution_file WHERE solution_id = $1", s.Id)
	saveSolutionFiles(ctx, tx, s)

	tx.Commit(ctx)
	return s
}