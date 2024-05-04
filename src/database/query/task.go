package query

import (
	"context"
	db "goto/src/database"
	"goto/src/model"
	"goto/src/utils"

	"github.com/jackc/pgx/v5"
)

func readTaskRow(row Scanable) *model.Task {
	task := model.Task{}
	task.InjectFiles = map[string]string{}
	err := row.Scan(
		&task.Id,
		&task.ProjectId,
		&task.Name,
		&task.Description,
		&task.RunTarget,
	)
	if err != nil {
		return nil
	}
	return &task
}

func readTaskRows(rows pgx.Rows) map[int]model.Task {
	tasksByIds := map[int]model.Task{}
	for rows.Next() {
		task := readTaskRow(rows)
		tasksByIds[task.Id] = *task
	}
	return tasksByIds
}

func readTaskRowThenExtend(ctx context.Context, row pgx.Row) *model.Task {
	task := readTaskRow(row)
	if task == nil {
		return nil
	}

	injectFiles := getInjectFilesByTasks(ctx, []int{task.Id})
	for _, ifl := range injectFiles {
		task.InjectFiles[ifl.Name] = ifl.Path
	}
	return task
}

func readTaskRowsThenExtend(ctx context.Context, rows pgx.Rows) []model.Task {
	tasksByIds := readTaskRows(rows)
	allInjectFiles := getInjectFilesByTasks(ctx, utils.MapKeys(tasksByIds))
	for _, ifl := range allInjectFiles {
		tasksByIds[ifl.TaskId].InjectFiles[ifl.Name] = ifl.Path
	}
	return utils.MapValues(tasksByIds)
}

func getTasksByProjects(ctx context.Context, projectIds []int) []model.Task {
	rows, _ := db.ConnPool.Query(
		ctx,
		"SELECT * FROM task WHERE project_id = ANY ($1)",
		projectIds,
	)
	return readTaskRowsThenExtend(ctx, rows)
}

func GetTask(ctx context.Context, id int) *model.Task {
	row := db.ConnPool.QueryRow(ctx, "SELECT * FROM task WHERE id = $1", id)
	return readTaskRowThenExtend(ctx, row)
}

func GetUserTasks(ctx context.Context, userId int) []model.Task {
	rows, _ := db.ConnPool.Query(ctx, "SELECT * FROM task WHERE user_id = $1", userId)
	return readTaskRowsThenExtend(ctx, rows)
}

func GetAllTasks(ctx context.Context) []model.Task {
	rows, _ := db.ConnPool.Query(ctx, "SELECT * FROM task")
	return readTaskRowsThenExtend(ctx, rows)
}

func getInjectFilesByTasks(ctx context.Context, taskIds []int) []model.InjectFile {
    allInjectFiles := []model.InjectFile{}

	rows, _ := db.ConnPool.Query(
		ctx,
		"SELECT id, task_id, name, path FROM injectfile WHERE task_id = ANY ($1)",
		taskIds,
	)

	for rows.Next() {
		injectFile := model.InjectFile{}
		rows.Scan(&injectFile.Id, &injectFile.TaskId, &injectFile.Name, &injectFile.Path)
		allInjectFiles = append(allInjectFiles, injectFile)
	}

	return allInjectFiles
}