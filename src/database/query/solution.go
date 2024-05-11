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
		&solution.UpdatedAt,
	)
	if err != nil {
		return nil
	}
	return &solution
}

func readSolutionRowWithResult(row Scanable) *model.Solution {
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
	solution := readSolutionRowWithResult(row)
	if solution == nil {
		return nil
	}

	solution.Files = getFilesBySolutionsWithCode(ctx, []int{solution.Id})
	return solution
}

func readSolutionRowsThenExtend(ctx context.Context, rows pgx.Rows) []model.Solution {
	solutionsByIds := readSolutionRows(rows)
	solutionFiles := getFilesBySolutions(ctx, utils.MapKeys(solutionsByIds))
	for _, sf := range solutionFiles {
		solution := solutionsByIds[sf.SolutionId]
		solution.Files = append(solution.Files, sf)
		solutionsByIds[sf.SolutionId] = solution
	}
	return utils.MapValues(solutionsByIds)
}

func GetSolution(ctx context.Context, id int) *model.Solution {
	row := db.ConnPool.QueryRow(ctx, "SELECT * FROM solution WHERE id = $1", id)
	return readSolutionRowThenExtend(ctx, row)
}

func GetUserSolution(ctx context.Context, id int, userId int) *model.Solution {
	row := db.ConnPool.QueryRow(
		ctx,
		"SELECT * FROM solution WHERE id = $1 AND user_id = $2",
		id, userId,
	)
	return readSolutionRowThenExtend(ctx, row)
}

func GetUserSolutions(ctx context.Context, userId int, pager *utils.Pager) []model.Solution {
	rows, _ := db.ConnPool.Query(
		ctx, `
		SELECT id, user_id, task_id, status, updated_at
        FROM solution WHERE user_id = $1`+pager.QuerySuffix(),
		userId,
	)
	return utils.MapValues(readSolutionRows(rows))
}

func getFilesBySolutions(ctx context.Context, solutionIds []int) []model.SolutionFile {
	solutionFiles := []model.SolutionFile{}

	rows, _ := db.ConnPool.Query(
		ctx, `
		SELECT id, solution_id, name
        FROM solution_file WHERE solution_id = ANY ($1)`,
		solutionIds,
	)

	for rows.Next() {
		sf := model.SolutionFile{}
		rows.Scan(&sf.Id, &sf.SolutionId, &sf.Name)
		solutionFiles = append(solutionFiles, sf)
	}

	return solutionFiles
}

func getFilesBySolutionsWithCode(ctx context.Context, solutionIds []int) []model.SolutionFile {
	solutionFiles := []model.SolutionFile{}

	rows, _ := db.ConnPool.Query(
		ctx, `
		SELECT id, solution_id, name, code
        FROM solution_file WHERE solution_id = ANY ($1)`,
		solutionIds,
	)

	for rows.Next() {
		sf := model.SolutionFile{}
		rows.Scan(&sf.Id, &sf.SolutionId, &sf.Name, &sf.Code)
		solutionFiles = append(solutionFiles, sf)
	}

	return solutionFiles
}

func saveSolutionFiles(ctx context.Context, tx pgx.Tx, s *model.Solution) {
	solutionFileEntries := make([][]any, len(s.Files))
	for i, sf := range s.Files {
		solutionFileEntries[i] = []any{s.Id, sf.Name, sf.Code}
	}
	tx.CopyFrom(
		ctx,
		pgx.Identifier{"solution_file"},
		[]string{"solution_id", "name", "code"},
		pgx.CopyFromRows(solutionFileEntries),
	)
}

func createSolution(ctx context.Context, tx pgx.Tx, s *model.Solution) {
	tx.QueryRow(
		ctx, `
        INSERT INTO solution (user_id, task_id, updated_at)
        VALUES ($1, $2, $3)
        RETURNING id`,
		s.UserId, s.TaskId, s.UpdatedAt,
	).Scan(&s.Id)
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
			createSolution(ctx, tx, s)
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