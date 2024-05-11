package query

import (
	"context"
	db "goto/src/database"
	m "goto/src/model"
	"goto/src/service"
	u "goto/src/utils"

	"github.com/jackc/pgx/v5"
)

func readProjectRow(row Scanable) *m.Project {
	project := m.Project{}
	err := row.Scan(
		&project.Id,
		&project.UserId,
		&project.Dir,
		&project.Name,
		&project.Language,
		&project.Containerization,
		&project.SrcDir,
		&project.StubDir,
	)
	if err != nil {
		return nil
	}
	return &project
}

func readProjectRows(rows pgx.Rows) map[int]m.Project {
	projectsByIds := map[int]m.Project{}
	for rows.Next() {
		project := readProjectRow(rows)
		projectsByIds[project.Id] = *project
	}
	return projectsByIds
}

func readProjectRowThenExtend(ctx context.Context, row pgx.Row) *m.Project {
	project := readProjectRow(row)
	if project == nil {
		return nil
	}
	project.Tasks = getTasksByProjectsWithStubs(ctx, []int{project.Id})
	project.Modules = getModulesByProjects(ctx, []int{project.Id}).Names()
	return project
}

func readProjectRowsThenExtend(ctx context.Context, rows pgx.Rows) m.Projects {
	projectsByIds := readProjectRows(rows)

	allTasks := getTasksByProjects(ctx, u.MapKeys(projectsByIds))
	for _, t := range allTasks {
		project := projectsByIds[t.ProjectId]
		project.Tasks = append(project.Tasks, t)
		projectsByIds[t.ProjectId] = project
	}

	allModules := getModulesByProjects(ctx, u.MapKeys(projectsByIds))
	for _, m := range allModules {
		project := projectsByIds[m.ProjectId]
		project.Modules = append(project.Modules, m.Name)
		projectsByIds[m.ProjectId] = project
	}

	return u.MapValues(projectsByIds)
}

func GetProject(ctx context.Context, id int) *m.Project {
	row := db.ConnPool.QueryRow(ctx, "SELECT * FROM project WHERE id = $1", id)
	return readProjectRowThenExtend(ctx, row)
}

func GetProjectShallow(ctx context.Context, id int) *m.Project {
	row := db.ConnPool.QueryRow(ctx, "SELECT * FROM project WHERE id = $1", id)
	return readProjectRow(row)
}

func GetUserProject(ctx context.Context, id int, userId int) *m.Project {
	row := db.ConnPool.QueryRow(
		ctx,
		"SELECT * FROM project WHERE id = $1 and user_id = $2",
		id, userId,
	)
	return readProjectRowThenExtend(ctx, row)
}

func GetUserProjects(ctx context.Context, userId int, pager *service.Pager) m.Projects {
	rows, _ := db.ConnPool.Query(
		ctx,
		"SELECT * FROM project WHERE user_id = $1"+pager.QuerySuffix(),
		userId,
	)
	return readProjectRowsThenExtend(ctx, rows)
}

func getModulesByProjects(ctx context.Context, projectIds []int) m.Modules {
	modules := m.Modules{}

	rows, _ := db.ConnPool.Query(
		ctx,
		"SELECT * FROM module WHERE project_id = ANY ($1)",
		projectIds,
	)

	for rows.Next() {
		m := m.Module{}
		rows.Scan(&m.Id, &m.ProjectId, &m.Name)
		modules = append(modules, m)
	}

	return modules
}

func CreateProject(ctx context.Context, p *m.Project) error { // TODO also return project
	tx, _ := db.ConnPool.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)

	projectEntries := [][]any{
		{p.UserId, p.Dir, p.Name, p.Language, p.Containerization, p.SrcDir, p.StubDir},
	}
	_, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"project"},
		[]string{"user_id", "dir", "name", "language", "containerization", "srcdir", "stubdir"},
		pgx.CopyFromRows(projectEntries),
	)
	if err != nil {
		return err
	}

	var projectId int
	err = tx.QueryRow(ctx, "SELECT id FROM project WHERE dir = $1", p.Dir).
		Scan(&projectId)
		// TODO remove after RETURNING id
	if err != nil {
		return err
	}

	moduleEntries := make([][]any, len(p.Modules))
	for i, mod := range p.Modules {
		moduleEntries[i] = []any{projectId, mod}
	}
	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"module"},
		[]string{"project_id", "name"},
		pgx.CopyFromRows(moduleEntries),
	)
	if err != nil {
		return err
	}

	taskEntries := make([][]any, len(p.Tasks))
	for i, t := range p.Tasks {
		taskEntries[i] = []any{projectId, t.Name, t.Description, t.RunTarget}
	}
	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"task"},
		[]string{"project_id", "name", "description", "runtarget"},
		pgx.CopyFromRows(taskEntries),
	)
	if err != nil {
		return err
	}

	taskFilesByTaskName := map[string]m.Task{}
	for _, t := range p.Tasks {
		taskFilesByTaskName[t.Name] = t
	}

	taskFileEntries := [][]any{}
	rows, err := tx.Query(ctx, "SELECT id, name FROM task WHERE project_id = $1", projectId)
	for rows.Next() {
		var taskId int
		var taskName string
		err = rows.Scan(&taskId, &taskName)
		if err != nil {
			return err
		}

		task := taskFilesByTaskName[taskName]
		for _, tf := range task.Files {
			taskFileEntries = append(taskFileEntries, []any{taskId, tf.Name, tf.Path, tf.Stub})
		}
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"task_file"},
		[]string{"task_id", "name", "path", "stub"},
		pgx.CopyFromRows(taskFileEntries),
	)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	return err
}

func DeleteProject(ctx context.Context, id int) {
	db.ConnPool.Exec(ctx, "DELETE FROM project WHERE id = $1", id)
}