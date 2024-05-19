package query

import (
	"context"
	"errors"
	db "goto/src/database"
	f "goto/src/filter"
	m "goto/src/model"
	"goto/src/service"
	u "goto/src/utils"
	"time"

	"github.com/jackc/pgx/v5"
)

const projectBaseSelectQuery = "SELECT * FROM project "

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
		&project.UpdatedAt,
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
	project.FailKeywords = GetFailKeywords(ctx, project.Id)
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

	allFailKeywords := getFailKeywordsByProjects(ctx, u.MapKeys(projectsByIds))
	for _, fk := range allFailKeywords {
		project := projectsByIds[fk.ProjectId]
		project.FailKeywords = append(project.FailKeywords, fk.Name)
		projectsByIds[fk.ProjectId] = project
	}

	return u.MapValues(projectsByIds)
}

func GetProjectShallow(ctx context.Context, id int) *m.Project {
	row := db.ConnPool.QueryRow(ctx, projectBaseSelectQuery+"WHERE id = $1", id)
	return readProjectRow(row)
}

func GetUserProject(ctx context.Context, id int, userId int) *m.Project {
	row := db.ConnPool.QueryRow(
		ctx,
		projectBaseSelectQuery+"WHERE id = $1 and user_id = $2",
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
		projectBaseSelectQuery+"WHERE"+filter.QueryCondition+pager.QuerySuffix,
		filter.QueryArgs...,
	)
	return readProjectRowsThenExtend(ctx, rows)
}

func GetFailKeywords(ctx context.Context, projectId int) []string {
	return getFailKeywordsByProjects(ctx, []int{projectId}).Names()
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

func getFailKeywordsByProjects(ctx context.Context, projectIds []int) m.FailKeywords {
	failKeywords := m.FailKeywords{}

	rows, _ := db.ConnPool.Query(
		ctx,
		"SELECT * FROM fail_keyword WHERE project_id = ANY ($1)",
		projectIds,
	)

	for rows.Next() {
		fk := m.FailKeyword{}
		rows.Scan(&fk.Id, &fk.ProjectId, &fk.Name)
		failKeywords = append(failKeywords, fk)
	}

	return failKeywords
}

func SaveProject(
	ctx context.Context,
	p *m.Project,
	cfg *m.GotoConfig,
) error {
	tx, _ := db.ConnPool.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)
	var err error

	if p.Id == 0 {
		err = createProjectOnly(ctx, tx, p)
	} else {
		err = updateProjectOnly(ctx, tx, p)
	}
	if err != nil {
		return err
	}

	err = saveProjectModules(ctx, tx, p.Id, p.Modules)
	if err != nil {
		return err
	}
	err = saveProjectFailKeywords(ctx, tx, p.Id, p.FailKeywords)
	if err != nil {
		return err
	}
	err = saveProjectTasks(ctx, tx, p.Id, p.Tasks, cfg.TaskConfigs)
	if err != nil {
		return err
	}
	err = saveProjectTaskFiles(ctx, tx, p.Id, p.Tasks)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	return err
}

func createProjectOnly(ctx context.Context, tx pgx.Tx, p *m.Project) error {
	err := tx.QueryRow(
		ctx, `
        INSERT INTO project
            (user_id, dir, name, language, containerization, srcdir, stubdir)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id`,
		p.UserId, p.Dir, p.Name, p.Language, p.Containerization, p.SrcDir, p.StubDir,
	).Scan(&p.Id)
	return err
}

func updateProjectOnly(ctx context.Context, tx pgx.Tx, p *m.Project) error {
	_, err := tx.Exec(
		ctx, `
        UPDATE project SET
            user_id = $1,
            dir = $2,
            name = $3,
            language = $4,
            containerization = $5,
            srcdir = $6,
            stubdir = $7,
            updated_at = $8
        WHERE id = $9`,
		p.UserId, p.Dir, p.Name, p.Language, p.Containerization, p.SrcDir, p.StubDir, time.Now(),
		p.Id,
	)
	return err
}

func saveProjectModules(ctx context.Context, tx pgx.Tx, projectId int, modules []string) error {
	tx.Exec(ctx, "DELETE FROM module WHERE project_id = $1", projectId)

	moduleEntries := make([][]any, len(modules))
	for i, m := range modules {
		moduleEntries[i] = []any{projectId, m}
	}

	_, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"module"},
		[]string{"project_id", "name"},
		pgx.CopyFromRows(moduleEntries),
	)

	return err
}

func saveProjectFailKeywords(
	ctx context.Context,
	tx pgx.Tx,
	projectId int,
	failKeywords []string,
) error {
	tx.Exec(ctx, "DELETE FROM fail_keyword WHERE project_id = $1", projectId)

	failKeywordEntries := make([][]any, len(failKeywords))
	for i, fk := range failKeywords {
		failKeywordEntries[i] = []any{projectId, fk}
	}

	_, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"fail_keyword"},
		[]string{"project_id", "name"},
		pgx.CopyFromRows(failKeywordEntries),
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
		if tc.OldName == "" {
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
		return errors.New(invalidTaskName + " task: `oldname` not found")
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
                    UNNEST($1::INTEGER[]) AS id,
                    UNNEST($2::VARCHAR(64)[]) AS name,
                    UNNEST($3::TEXT[]) AS description,
                    UNNEST($4::VARCHAR(256)[]) AS runtarget
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