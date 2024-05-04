package query

import (
	"context"
	db "goto/src/database"
	"goto/src/model"

	"github.com/jackc/pgx/v5"
)

func readSolutionRow(row Scanable) *model.Solution {
	solution := model.Solution{}
	err := row.Scan(
		&solution.Id,
		&solution.UserId,
		&solution.TaskId,
		&solution.Status,
		&solution.Code,
		&solution.Result,
		&solution.UpdatedAt,
	)
	if err != nil {
		return nil
	}
	return &solution
}

func readSolutionRows(rows pgx.Rows) []model.Solution {
    solutions := []model.Solution{}
	for rows.Next() {
		solutions = append(solutions, *readSolutionRow(rows))
	}
	return solutions
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
	return readSolutionRow(row)
}

func GetUserSolutions(ctx context.Context, userId int) []model.Solution {
	rows, _ := db.ConnPool.Query(ctx, "SELECT * FROM solution WHERE user_id = $1", userId)
	return readSolutionRows(rows)
}