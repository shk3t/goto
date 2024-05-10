package query

import (
	"context"
	db "goto/src/database"
	"goto/src/model"
	"goto/src/utils"

	"github.com/jackc/pgx/v5"
)

func readProjectRow(row Scanable) *model.Project {
	project := model.Project{}
	err := row.Scan(
		&project.Id,
		&project.User.Id,
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

func readProjectRows(rows pgx.Rows) map[int]model.Project {
	projectsByIds := map[int]model.Project{}
	for rows.Next() {
		project := readProjectRow(rows)
		projectsByIds[project.Id] = *project
	}
	return projectsByIds
}

func readProjectRowThenExtend(ctx context.Context, row pgx.Row) *model.Project {
	project := readProjectRow(row)
	if project == nil {
		return nil
	}

	project.Tasks = getTasksByProjectsWithStubs(ctx, []int{project.Id})
	modules := getModulesByProjects(ctx, []int{project.Id})
	project.Modules = make([]string, len(modules))
	for i, m := range modules {
		project.Modules[i] = m.Name
	}
	return project
}

func readProjectRowsThenExtend(ctx context.Context, rows pgx.Rows) []model.Project {
	projectsByIds := readProjectRows(rows)

	allTasks := getTasksByProjects(ctx, utils.MapKeys(projectsByIds))
	for _, t := range allTasks {
		project := projectsByIds[t.ProjectId]
		project.Tasks = append(project.Tasks, t)
		projectsByIds[t.ProjectId] = project
	}

	allModules := getModulesByProjects(ctx, utils.MapKeys(projectsByIds))
	for _, m := range allModules {
		project := projectsByIds[m.ProjectId]
		project.Modules = append(project.Modules, m.Name)
		projectsByIds[m.ProjectId] = project
	}

	return utils.MapValues(projectsByIds)
}

func GetProject(ctx context.Context, id int) *model.Project {
	row := db.ConnPool.QueryRow(ctx, "SELECT * FROM project WHERE id = $1", id)
	return readProjectRowThenExtend(ctx, row)
}

func GetProjectShallow(ctx context.Context, id int) *model.Project {
	row := db.ConnPool.QueryRow(ctx, "SELECT * FROM project WHERE id = $1", id)
	return readProjectRow(row)
}

func GetUserProject(ctx context.Context, id int, userId int) *model.Project {
	row := db.ConnPool.QueryRow(
		ctx,
		"SELECT * FROM project WHERE id = $1 and user_id = $2",
		id, userId,
	)
	return readProjectRowThenExtend(ctx, row)
}

func GetUserProjects(ctx context.Context, userId int, pager *utils.Pager) []model.Project {
	rows, _ := db.ConnPool.Query(
		ctx,
		"SELECT * FROM project WHERE user_id = $1"+pager.QuerySuffix(),
		userId,
	)
	return readProjectRowsThenExtend(ctx, rows)
}

func getModulesByProjects(ctx context.Context, projectIds []int) []model.Module {
	modules := []model.Module{}

	rows, _ := db.ConnPool.Query(
		ctx,
		"SELECT * FROM module WHERE project_id = ANY ($1)",
		projectIds,
	)

	for rows.Next() {
		m := model.Module{}
		rows.Scan(&m.Id, &m.ProjectId, &m.Name)
		modules = append(modules, m)
	}

	return modules
}

func CreateProject(ctx context.Context, p *model.Project) error { // TODO also return project
	tx, _ := db.ConnPool.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)

	projectEntries := [][]any{
		{p.User.Id, p.Dir, p.Name, p.Language, p.Containerization, p.SrcDir, p.StubDir},
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

	taskFilesByTaskName := map[string]model.Task{}
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