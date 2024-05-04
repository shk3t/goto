package query

import (
	"context"
	db "goto/src/database"
	"goto/src/model"
	"goto/src/utils"

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