package query

import (
	"context"
	db "goto/src/database"
	"goto/src/model"
	"goto/src/utils"

	"github.com/jackc/pgx/v5"
)

type Scanable interface {
	Scan(dest ...any) error
}

func readTaskRow(row Scanable) (*model.Task, error) {
	task := model.Task{}
	task.InjectFiles = make(map[string]string)
	err := row.Scan(
		&task.Id,
		&task.ProjectId,
		&task.Name,
		&task.Description,
		&task.RunTarget,
	)
	return &task, err
}

func readTaskRows(rows pgx.Rows) map[int]model.Task {
	tasksByIds := make(map[int]model.Task)
	for rows.Next() {
		task, _ := readTaskRow(rows)
		tasksByIds[task.Id] = *task
	}
	return tasksByIds
}

func getTasksByProjects(ctx context.Context, projectIds []int) []model.Task {
	rows, _ := db.ConnPool.Query(
		ctx,
		"SELECT * FROM task WHERE project_id = ANY ($1)",
		projectIds,
	)

	tasksByIds := readTaskRows(rows)
	extendTasksWithInjectFiles(ctx, tasksByIds)

	return utils.MapValues(tasksByIds)
}

func extendTaskWithInjectFiles(ctx context.Context, task *model.Task) {
	injectFiles := getInjectFilesByTasks(ctx, []int{task.Id})
	for _, ifl := range injectFiles {
		task.InjectFiles[ifl.Name] = ifl.Path
	}
}

func extendTasksWithInjectFiles(ctx context.Context, tasksByIds map[int]model.Task) {
	allInjectFiles := getInjectFilesByTasks(ctx, utils.MapKeys(tasksByIds))
	for _, ifl := range allInjectFiles {
		tasksByIds[ifl.TaskId].InjectFiles[ifl.Name] = ifl.Path
	}
}

func getInjectFilesByTasks(ctx context.Context, taskIds []int) []model.InjectFile {
	var allInjectFiles []model.InjectFile

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

func GetTask(ctx context.Context, id int) (*model.Task, error) {
	row := db.ConnPool.QueryRow(
		ctx, "SELECT * FROM task WHERE id = $1", id,
	)
	task, err := readTaskRow(row)
	if err != nil {
		return nil, err
	}
	extendTaskWithInjectFiles(ctx, task)
	return task, nil
}

func GetUserTasks(ctx context.Context, userId int) []model.Task {
	rows, _ := db.ConnPool.Query(ctx, "SELECT * FROM task WHERE user_id = $1", userId)
	tasksByIds := readTaskRows(rows)
	extendTasksWithInjectFiles(ctx, tasksByIds)
	return utils.MapValues(tasksByIds)
}

func GetAllTasks(ctx context.Context) []model.Task {
	rows, _ := db.ConnPool.Query(ctx, "SELECT * FROM task")
	tasksByIds := readTaskRows(rows)
	extendTasksWithInjectFiles(ctx, tasksByIds)
	return utils.MapValues(tasksByIds)
}