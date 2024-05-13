package query

import (
	"context"
	db "goto/src/database"
	f "goto/src/filter"
	m "goto/src/model"
	"goto/src/service"
	u "goto/src/utils"
	"time"

	"github.com/jackc/pgx/v5"
)

func readSolutionRowBase(row Scanable, withResult bool) *m.Solution {
	solution := m.Solution{}
	args := []any{
		&solution.Id,
		&solution.UserId,
		&solution.Status,
		&solution.UpdatedAt,
		&solution.Task.Id,
		&solution.Task.ProjectId,
		&solution.Task.Name,
		&solution.Task.Description,
		nil,
		&solution.Task.Language,
		&solution.Task.UpdatedAt,
	}
	if withResult {
		args = u.Insert(args, 3, &solution.Result)
	}
	if err := row.Scan(args...); err != nil {
		return nil
	}
	return &solution
}

func readSolutionRow(row Scanable) *m.Solution {
	return readSolutionRowBase(row, false)
}

func readSolutionRowWithResult(row Scanable) *m.Solution {
	return readSolutionRowBase(row, true)
}

func readSolutionRows(rows pgx.Rows) map[int]m.Solution {
	solutionsByIds := map[int]m.Solution{}
	for rows.Next() {
		solution := readSolutionRow(rows)
		solutionsByIds[solution.Id] = *solution
	}
	return solutionsByIds
}

func readSolutionRowThenExtend(ctx context.Context, row pgx.Row) *m.Solution {
	solution := readSolutionRowWithResult(row)
	if solution == nil {
		return nil
	}

	solution.Files = getFilesBySolutionsWithCode(ctx, []int{solution.Id})
	solution.Task.FileNames = getFilesByTasks(ctx, []int{solution.Task.Id}).Names()
	solution.Task.Modules = getModulesByTasks(ctx, []int{solution.Task.Id}).Names()
	return solution
}

func readSolutionRowsThenExtend(ctx context.Context, rows pgx.Rows) m.Solutions {
	solutionsByIds := readSolutionRows(rows)

	solutionsByTaskIds := map[int]m.Solution{}
	for _, s := range solutionsByIds {
		solutionsByTaskIds[s.Task.Id] = s
	}

	allTaskFiles := getFilesByTasks(ctx, u.MapKeys(solutionsByTaskIds))
	for _, tf := range allTaskFiles {
		solution := solutionsByTaskIds[tf.TaskId]
		solution.Task.FileNames = append(solution.Task.FileNames, tf.Name)
		solutionsByTaskIds[tf.TaskId] = solution
	}

	allModules := getModulesByTasks(ctx, u.MapKeys(solutionsByTaskIds))
	for _, m := range allModules {
		solution := solutionsByTaskIds[m.TaskId]
		solution.Task.Modules = append(solution.Task.Modules, m.Name)
		solutionsByTaskIds[m.TaskId] = solution
	}

	return u.MapValues(solutionsByTaskIds)
}

func GetUserSolution(
	ctx context.Context,
	id int,
	userId int,
) *m.Solution {
	row := db.ConnPool.QueryRow(
		ctx, `
        SELECT
            solution.id,
            solution.user_id,
            solution.status,
            solution.result,
            solution.updated_at,
            task.*,
            project.language,
            project.updated_at
        FROM solution
        JOIN task ON task.id = solution.task_id
        JOIN project ON project.id = task.project_id
        WHERE solution.id = $1 AND solution.user_id = $2`,
		id, userId,
	)
	return readSolutionRowThenExtend(ctx, row)
}

func GetSolutions(
	ctx context.Context,
	userId int,
	pager *service.Pager,
	filter *f.SolutionFilter,
) m.Solutions {
	rows, _ := db.ConnPool.Query(
		ctx, `
		SELECT
            solution.id,
            solution.user_id,
            solution.status,
            solution.updated_at,
            task.*,
            project.language,
            project.updated_at
        FROM solution
        JOIN task ON task.id = solution.task_id
        JOIN project ON project.id = task.project_id
        WHERE`+filter.SqlCondition+pager.QuerySuffix,
		filter.SqlArgs...,
	)
	return readSolutionRowsThenExtend(ctx, rows)
}

func getFilesBySolutions(ctx context.Context, solutionIds []int) m.SolutionFiles {
	solutionFiles := m.SolutionFiles{}

	rows, _ := db.ConnPool.Query(
		ctx, `
		SELECT id, solution_id, name
        FROM solution_file WHERE solution_id = ANY ($1)`,
		solutionIds,
	)

	for rows.Next() {
		sf := m.SolutionFile{}
		rows.Scan(&sf.Id, &sf.SolutionId, &sf.Name)
		solutionFiles = append(solutionFiles, sf)
	}

	return solutionFiles
}

func getFilesBySolutionsWithCode(ctx context.Context, solutionIds []int) m.SolutionFiles {
	solutionFiles := m.SolutionFiles{}

	rows, _ := db.ConnPool.Query(
		ctx, `
		SELECT id, solution_id, name, code
        FROM solution_file WHERE solution_id = ANY ($1)`,
		solutionIds,
	)

	for rows.Next() {
		sf := m.SolutionFile{}
		rows.Scan(&sf.Id, &sf.SolutionId, &sf.Name, &sf.Code)
		solutionFiles = append(solutionFiles, sf)
	}

	return solutionFiles
}

func saveSolutionFiles(ctx context.Context, tx pgx.Tx, s *m.Solution) {
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

func createSolution(ctx context.Context, tx pgx.Tx, s *m.Solution) {
	tx.QueryRow(
		ctx, `
        INSERT INTO solution (user_id, task_id, updated_at)
        VALUES ($1, $2, $3)
        RETURNING id`,
		s.UserId, s.Task.Id, s.UpdatedAt,
	).Scan(&s.Id)
}

func updateSolution(ctx context.Context, tx pgx.Tx, s *m.Solution) {
	tx.Exec(
		ctx, `
        UPDATE solution
        SET status = $1, result = $2, updated_at = $3
        WHERE id = $4`,
		s.Status, s.Result, s.UpdatedAt,
		s.Id,
	)
}

func SaveSolution(ctx context.Context, s *m.Solution) {
	tx, _ := db.ConnPool.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)
	s.UpdatedAt = time.Now()

	if s.Id == 0 {
		s.Status = "check"

		tx.QueryRow(
			ctx,
			"SELECT id from solution WHERE user_id = $1 AND task_id = $2",
			s.UserId, s.Task.Id,
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
}