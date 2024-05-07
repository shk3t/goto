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
	task.Files = map[string]string{}
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

	files := getFilesByTasks(ctx, []int{task.Id})
	for _, tf := range files {
		task.Files[tf.Name] = tf.Path
	}
	return task
}

func readTaskRowsThenExtend(ctx context.Context, rows pgx.Rows) []model.Task {
	tasksByIds := readTaskRows(rows)
	files := getFilesByTasks(ctx, utils.MapKeys(tasksByIds))
	for _, tf := range files {
		tasksByIds[tf.TaskId].Files[tf.Name] = tf.Path
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

func GetTaskWithStubs(ctx context.Context, id int) *model.Task {
	row := db.ConnPool.QueryRow(ctx, "SELECT * FROM task WHERE id = $1", id)
	return readTaskRowThenExtend(ctx, row)
}

func GetUserTasks(ctx context.Context, userId int, pager *utils.Pager) []model.Task {
	rows, _ := db.ConnPool.Query(
		ctx,
		"SELECT * FROM task WHERE user_id = $1"+pager.QuerySuffix(),
		userId,
	)
	return readTaskRowsThenExtend(ctx, rows)
}

func GetAllTasks(ctx context.Context, pager *utils.Pager) []model.Task {
	rows, _ := db.ConnPool.Query(ctx, "SELECT * FROM task"+pager.QuerySuffix())
	return readTaskRowsThenExtend(ctx, rows)
}

func getFilesByTasks(ctx context.Context, taskIds []int) []model.TaskFile {
	files := []model.TaskFile{}

	rows, _ := db.ConnPool.Query(
		ctx,
		"SELECT id, task_id, name, path FROM task_file WHERE task_id = ANY ($1)",
		taskIds,
	)

	for rows.Next() {
		file := model.TaskFile{}
		rows.Scan(&file.Id, &file.TaskId, &file.Name, &file.Path)
		files = append(files, file)
	}

	return files
}