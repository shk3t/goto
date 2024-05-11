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
	err := row.Scan(
		&task.Id,
		&task.ProjectId,
		&task.Name,
		&task.Description,
		&task.RunTarget,
		&task.Language,
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
	task.Files = getFilesByTasksWithStubs(ctx, []int{task.Id})
	modules := getModulesByTasks(ctx, []int{task.Id})
	task.Modules = make([]string, len(modules))
	for i, m := range modules {
		task.Modules[i] = m.Name
	}
	return task
}

func readTaskRowsThenExtendBase(ctx context.Context, rows pgx.Rows, withStubs bool) []model.Task {
	tasksByIds := readTaskRows(rows)

	var taskFiles []model.TaskFile
	if withStubs {
		taskFiles = getFilesByTasksWithStubs(ctx, utils.MapKeys(tasksByIds))
	} else {
		taskFiles = getFilesByTasks(ctx, utils.MapKeys(tasksByIds))
	}

	for _, tf := range taskFiles {
		task := tasksByIds[tf.TaskId]
		task.Files = append(task.Files, tf)
		tasksByIds[tf.TaskId] = task
	}

	allModules := getModulesByTasks(ctx, utils.MapKeys(tasksByIds))
	for _, m := range allModules {
		task := tasksByIds[m.TaskId]
		task.Modules = append(task.Modules, m.Name)
		tasksByIds[m.TaskId] = task
	}

	return utils.MapValues(tasksByIds)
}

func readTaskRowsThenExtend(ctx context.Context, rows pgx.Rows) []model.Task {
	return readTaskRowsThenExtendBase(ctx, rows, false)
}

func readTaskRowsThenExtendWithStubs(ctx context.Context, rows pgx.Rows) []model.Task {
	return readTaskRowsThenExtendBase(ctx, rows, true)
}

func GetTask(ctx context.Context, id int) *model.Task {
	row := db.ConnPool.QueryRow(
		ctx, `
        SELECT t.*, p.language
        FROM task as t
        JOIN project AS p ON p.id = t.project_id
        WHERE t.id = $1`,
		id,
	)
	return readTaskRowThenExtend(ctx, row)
}

func GetUserTasks(ctx context.Context, userId int, pager *utils.Pager) []model.Task {
	rows, _ := db.ConnPool.Query(
		ctx, `
        SELECT t.*, p.language
        FROM task as t
        JOIN project AS p ON p.id = t.project_id
        WHERE p.user_id = $1`+
			pager.QuerySuffix(),
		userId,
	)
	return readTaskRowsThenExtend(ctx, rows)
}

func GetAllTasks(ctx context.Context, pager *utils.Pager) []model.Task {
	rows, _ := db.ConnPool.Query(
		ctx, `
        SELECT t.*, p.language
        FROM task as t
        JOIN project AS p ON p.id = t.project_id`+
			pager.QuerySuffix(),
	)
	return readTaskRowsThenExtend(ctx, rows)
}

func getTasksByProjects(ctx context.Context, projectIds []int) []model.Task {
	rows, _ := db.ConnPool.Query(
		ctx, `
        SELECT t.*, p.language
        FROM task as t
        JOIN project AS p ON p.id = t.project_id
        WHERE p.id = ANY ($1)`,
		projectIds,
	)
	return readTaskRowsThenExtend(ctx, rows)
}

func getTasksByProjectsWithStubs(ctx context.Context, projectIds []int) []model.Task {
	rows, _ := db.ConnPool.Query(
		ctx, `
        SELECT t.*, p.language
        FROM task as t
        JOIN project AS p ON p.id = t.project_id
        WHERE p.id = ANY ($1)`,
		projectIds,
	)
	return readTaskRowsThenExtendWithStubs(ctx, rows)
}

func getFilesByTasks(ctx context.Context, taskIds []int) []model.TaskFile {
	taskFiles := []model.TaskFile{}

	rows, _ := db.ConnPool.Query(
		ctx, `
		SELECT id, task_id, name, path
        FROM task_file WHERE task_id = ANY ($1)`,
		taskIds,
	)

	for rows.Next() {
		tf := model.TaskFile{}
		rows.Scan(&tf.Id, &tf.TaskId, &tf.Name, &tf.Path)
		taskFiles = append(taskFiles, tf)
	}

	return taskFiles
}

func getFilesByTasksWithStubs(ctx context.Context, taskIds []int) []model.TaskFile {
	taskFiles := []model.TaskFile{}

	rows, _ := db.ConnPool.Query(
		ctx, `
		SELECT id, task_id, name, path, stub
        FROM task_file WHERE task_id = ANY ($1)`,
		taskIds,
	)

	for rows.Next() {
		tf := model.TaskFile{}
		rows.Scan(&tf.Id, &tf.TaskId, &tf.Name, &tf.Path, &tf.Stub)
		taskFiles = append(taskFiles, tf)
	}

	return taskFiles
}

func getModulesByTasks(ctx context.Context, taskIds []int) []model.Module {
	modules := []model.Module{}

	rows, _ := db.ConnPool.Query(
		ctx, `
		SELECT m.id, t.id, m.name
        FROM module AS m
        JOIN project AS p ON p.id = m.project_id
        JOIN task AS t ON t.project_id = p.id
        WHERE t.id = ANY ($1)`,
		taskIds,
	)

	for rows.Next() {
		m := model.Module{}
		rows.Scan(&m.Id, &m.TaskId, &m.Name)
		modules = append(modules, m)
	}

	return modules
}