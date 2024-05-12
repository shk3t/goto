package query

import (
	"context"
	"errors"
	db "goto/src/database"
	f "goto/src/filter"
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

func GetProjects(
	ctx context.Context,
	pager *service.Pager,
	filter *f.ProjectFilter,
) m.Projects {
	rows, _ := db.ConnPool.Query(
		ctx,
		"SELECT * FROM project WHERE"+filter.SqlCondition+pager.QuerySuffix,
		filter.SqlArgs...,
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

func CreateProject(
	ctx context.Context,
	p *m.Project,
	cfg *m.GotoConfig,
) (*m.Project, error) {
	tx, _ := db.ConnPool.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)

	err := db.ConnPool.QueryRow(
		ctx, `
        INSERT INTO project
            (user_id, dir, name, language, containerization, srcdir, stubdir)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id`,
		p.UserId, p.Dir, p.Name, p.Language, p.Containerization, p.SrcDir, p.StubDir,
	).Scan(&p.Id)
	if err != nil {
		return nil, err
	}

	saveProjectModules(ctx, tx, p.Id, p.Modules)
	saveProjectTasks(ctx, tx, p.Id, p.Tasks, cfg.TaskConfigs)
	saveProjectTaskFiles(ctx, tx, p.Id, p.Tasks)

	err = tx.Commit(ctx)
	return p, err
}

func UpdateProject(
	ctx context.Context,
	p *m.Project,
	cfg *m.GotoConfig,
) (*m.Project, error) {
	tx, _ := db.ConnPool.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)

	_, err := db.ConnPool.Exec(
		ctx, `
        UPDATE project SET
            user_id = $1,
            dir = $2,
            name = $3,
            language = $4,
            containerization = $6,
            srcdir = $7,
            stubdir = $8
        WHERE id = $9`,
		p.UserId, p.Dir, p.Name, p.Language, p.Containerization, p.SrcDir, p.StubDir,
		p.Id,
	)
	if err != nil {
		return nil, err
	}

	saveProjectModules(ctx, tx, p.Id, p.Modules)
	saveProjectTasks(ctx, tx, p.Id, p.Tasks, cfg.TaskConfigs)
	saveProjectTaskFiles(ctx, tx, p.Id, p.Tasks)

	err = tx.Commit(ctx)
	return p, err
}

func saveProjectModules(ctx context.Context, tx pgx.Tx, projectId int, modules []string) error {
	tx.Exec(ctx, "DELETE FROM module WHERE project_id = $1", projectId)

	moduleEntries := make([][]any, len(modules))
	for i, mod := range modules {
		moduleEntries[i] = []any{projectId, mod}
	}

	_, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"module"},
		[]string{"project_id", "name"},
		pgx.CopyFromRows(moduleEntries),
	)

	return err
}

func saveProjectTasks(
	ctx context.Context,
	tx pgx.Tx,
	projectId int,
	tasks m.Tasks,
	taskConfigs m.TaskConfigs,
) error {
	existingTaskIdsByNames := map[string]int{}
	rows, err := tx.Query(ctx, "SELECT id, name FROM task WHERE project_id = $1", projectId)
	for rows.Next() {
		var taskId int
		var taskName string
		err = rows.Scan(&taskId, &taskName)
		if err != nil {
			return err
		}
		existingTaskIdsByNames[taskName] = taskId
	}
	existingTaskNames := u.MapKeys(existingTaskIdsByNames)

	specifiedTasksByNames := map[string]m.Task{}
	for _, t := range tasks {
		specifiedTasksByNames[t.Name] = t
	}
	specifiedTasksByOldNames := map[string]m.Task{}
	specifiedOldNames := []string{}
	for _, tc := range taskConfigs {
		if tc.OldName != "" {
			specifiedTasksByOldNames[tc.Name] = specifiedTasksByNames[tc.Name]
		} else {
			specifiedTasksByOldNames[tc.OldName] = specifiedTasksByNames[tc.Name]
			specifiedOldNames = append(specifiedOldNames, tc.OldName)
		}
	}
	taskOldNames := u.MapKeys(specifiedTasksByOldNames)

	invalidOldNames := u.Difference(specifiedOldNames, existingTaskNames)
	if len(invalidOldNames) > 0 {
		invalidTaskName := specifiedTasksByOldNames[invalidOldNames[0]].Name
		return errors.New("Task" + invalidTaskName + ": `oldname` not found")
	}

	deprecatedTaskNames := u.Difference(existingTaskNames, taskOldNames)
	if len(deprecatedTaskNames) > 0 {
		tx.Exec(
			ctx,
			"DELETE FROM task WHERE project_id = $1 AND name = ANY ($2)",
			projectId, deprecatedTaskNames,
		)
	}

	updatedTaskNames := u.Intersection(existingTaskNames, taskOldNames)
	if len(updatedTaskNames) > 0 {
		ids := make([]int, len(updatedTaskNames))
		names := make([]string, len(updatedTaskNames))
		descriptions := make([]string, len(updatedTaskNames))
		runtargets := make([]string, len(updatedTaskNames))
		for i, name := range updatedTaskNames {
			task := specifiedTasksByOldNames[name]
			ids[i] = existingTaskIdsByNames[name]
			names[i] = task.Name
			descriptions[i] = task.Description
			runtargets[i] = task.RunTarget
		}
		_, err := tx.Exec(
			ctx, `
            UPDATE task SET
                name = data.name,
                description = data.description,
                runtarget = data.runtarget
            FROM (
                SELECT
                    UNNEST($1) AS id,
                    UNNEST($2) AS name,
                    UNNEST($3) AS description,
                    UNNEST($4) AS runtarget,
            ) AS data
            WHERE task.id = data.id`,
			ids, names, descriptions, runtargets,
		)
		if err != nil {
			return err
		}
	}

	newTaskNames := u.Difference(taskOldNames, existingTaskNames)
	if len(newTaskNames) > 0 {
		taskEntries := make([][]any, len(newTaskNames))
		for i, name := range newTaskNames {
			t := specifiedTasksByOldNames[name]
			taskEntries[i] = []any{projectId, t.Name, t.Description, t.RunTarget}
		}
		_, err = tx.CopyFrom(
			ctx,
			pgx.Identifier{"task"},
			[]string{"project_id", "name", "description", "runtarget"},
			pgx.CopyFromRows(taskEntries),
		)
	}
	return err
}

func saveProjectTaskFiles(ctx context.Context, tx pgx.Tx, projectId int, tasks m.Tasks) error {
	tx.Exec(
		ctx, `
        DELETE FROM task_file
        USING task
        WHERE task.id = task_file.task_id
            AND task.project_id = $1`,
		projectId,
	)

	taskFilesByTaskName := map[string]m.Task{}
	for _, t := range tasks {
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
	return err
}

func DeleteProject(ctx context.Context, id int) {
	db.ConnPool.Exec(ctx, "DELETE FROM project WHERE id = $1", id)
}