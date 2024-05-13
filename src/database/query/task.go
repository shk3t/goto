package query

import (
	"context"
	db "goto/src/database"
	f "goto/src/filter"
	m "goto/src/model"
	"goto/src/service"
	u "goto/src/utils"

	"github.com/jackc/pgx/v5"
)

func readTaskRow(row Scanable) *m.Task {
	task := m.Task{}
	err := row.Scan(
		&task.Id,
		&task.ProjectId,
		&task.Name,
		&task.Description,
		&task.RunTarget,
		&task.Language,
		&task.UpdatedAt,
	)
	if err != nil {
		return nil
	}
	return &task
}

func readTaskRows(rows pgx.Rows) map[int]m.Task {
	tasksByIds := map[int]m.Task{}
	for rows.Next() {
		task := readTaskRow(rows)
		tasksByIds[task.Id] = *task
	}
	return tasksByIds
}

func readTaskRowThenExtend(ctx context.Context, row pgx.Row) *m.Task {
	task := readTaskRow(row)
	if task == nil {
		return nil
	}
	task.Files = getFilesByTasksWithStubs(ctx, []int{task.Id})
	task.Modules = getModulesByTasks(ctx, []int{task.Id}).Names()
	return task
}

func readTaskRowsThenExtendBase(ctx context.Context, rows pgx.Rows, withStubs bool) m.Tasks {
	tasksByIds := readTaskRows(rows)

	var allTaskFiles m.TaskFiles
	if withStubs {
		allTaskFiles = getFilesByTasksWithStubs(ctx, u.MapKeys(tasksByIds))
	} else {
		allTaskFiles = getFilesByTasks(ctx, u.MapKeys(tasksByIds))
	}

	for _, tf := range allTaskFiles {
		task := tasksByIds[tf.TaskId]
		task.Files = append(task.Files, tf)
		tasksByIds[tf.TaskId] = task
	}

	allModules := getModulesByTasks(ctx, u.MapKeys(tasksByIds))
	for _, m := range allModules {
		task := tasksByIds[m.TaskId]
		task.Modules = append(task.Modules, m.Name)
		tasksByIds[m.TaskId] = task
	}

	return u.MapValues(tasksByIds)
}

func readTaskRowsThenExtend(ctx context.Context, rows pgx.Rows) m.Tasks {
	return readTaskRowsThenExtendBase(ctx, rows, false)
}

func readTaskRowsThenExtendWithStubs(ctx context.Context, rows pgx.Rows) m.Tasks {
	return readTaskRowsThenExtendBase(ctx, rows, true)
}

func GetTask(ctx context.Context, id int) *m.Task {
	row := db.ConnPool.QueryRow(
		ctx, `
        SELECT task.*, project.language, project.updated_at
        FROM task
        JOIN project ON project.id = task.project_id
        WHERE task.id = $1`,
		id,
	)
	return readTaskRowThenExtend(ctx, row)
}

func GetTasks(ctx context.Context, pager *service.Pager, filter *f.TaskFilter) m.Tasks {
	rows, _ := db.ConnPool.Query(
		ctx, `
        SELECT task.*, project.language, project.updated_at
        FROM task
        JOIN project ON project.id = task.project_id
        WHERE`+filter.SqlCondition+
			pager.QuerySuffix,
		filter.SqlArgs...,
	)
	return readTaskRowsThenExtend(ctx, rows)
}

func getTasksByProjects(ctx context.Context, projectIds []int) m.Tasks {
	rows, _ := db.ConnPool.Query(
		ctx, `
        SELECT task.*, project.language, project.updated_at
        FROM task
        JOIN project ON project.id = task.project_id
        WHERE project.id = ANY ($1)`,
		projectIds,
	)
	return readTaskRowsThenExtend(ctx, rows)
}

func getTasksByProjectsWithStubs(ctx context.Context, projectIds []int) m.Tasks {
	rows, _ := db.ConnPool.Query(
		ctx, `
        SELECT task.*, project.language, project.updated_at
        FROM task
        JOIN project ON project.id = task.project_id
        WHERE project.id = ANY ($1)`,
		projectIds,
	)
	return readTaskRowsThenExtendWithStubs(ctx, rows)
}

func getFilesByTasks(ctx context.Context, taskIds []int) m.TaskFiles {
	taskFiles := m.TaskFiles{}

	rows, _ := db.ConnPool.Query(
		ctx, `
		SELECT id, task_id, name, path
        FROM task_file WHERE task_id = ANY ($1)`,
		taskIds,
	)

	for rows.Next() {
		tf := m.TaskFile{}
		rows.Scan(&tf.Id, &tf.TaskId, &tf.Name, &tf.Path)
		taskFiles = append(taskFiles, tf)
	}

	return taskFiles
}

func getFilesByTasksWithStubs(ctx context.Context, taskIds []int) m.TaskFiles {
	taskFiles := m.TaskFiles{}

	rows, _ := db.ConnPool.Query(
		ctx, `
		SELECT id, task_id, name, path, stub
        FROM task_file WHERE task_id = ANY ($1)`,
		taskIds,
	)

	for rows.Next() {
		tf := m.TaskFile{}
		rows.Scan(&tf.Id, &tf.TaskId, &tf.Name, &tf.Path, &tf.Stub)
		taskFiles = append(taskFiles, tf)
	}

	return taskFiles
}

func getModulesByTasks(ctx context.Context, taskIds []int) m.Modules {
	modules := m.Modules{}

	rows, _ := db.ConnPool.Query(
		ctx, `
		SELECT module.id, task.id, module.name
        FROM module
        JOIN project ON project.id = module.project_id
        JOIN task ON task.project_id = project.id
        WHERE task.id = ANY ($1)`,
		taskIds,
	)

	for rows.Next() {
		m := m.Module{}
		rows.Scan(&m.Id, &m.TaskId, &m.Name)
		modules = append(modules, m)
	}

	return modules
}