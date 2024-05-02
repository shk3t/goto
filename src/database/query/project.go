package query

import (
	"context"
	"errors"
	db "goto/src/database"
	"goto/src/model"
	"goto/src/utils"

	"github.com/jackc/pgx/v5"
)

func GetProject(ctx context.Context, id int) (*model.Project, error) {
	projects, _ := getProjects(ctx, 0, []int{id})
	if len(projects) == 0 {
		return nil, errors.New("Not found")
	}
	return &projects[0], nil
}

func GetUserProject(ctx context.Context, userId int, id int) (*model.Project, error) {
	projects, _ := getProjects(ctx, userId, []int{id})
	if len(projects) == 0 {
		return nil, errors.New("Not found")
	}
	return &projects[0], nil
}

func GetUserProjects(ctx context.Context, userId int) ([]model.Project, error) {
	projects, err := getProjects(ctx, userId, nil)
	return projects, err
}

func getProjects(ctx context.Context, userId int, projectIds []int) ([]model.Project, error) {
	projectsByIds := make(map[int]model.Project)

	query := "SELECT * FROM project WHERE "
	var params []any
	if userId != 0 {
		if projectIds != nil {
			query += "user_id = $1 AND id = ANY ($2)"
			params = []any{userId, projectIds}
		} else {
			query += "user_id = $1"
			params = []any{userId}
		}
	} else {
		if projectIds != nil {
			query += "id = ANY ($2)"
			params = []any{projectIds}
		} else {
			return nil, errors.New("Not supported")
		}
	}

	rows, _ := db.ConnPool.Query(ctx, query, params...)

	for rows.Next() {
		project := model.Project{}
		rows.Scan(
			&project.Id,
			&project.User.Id,
			&project.Dir,
			&project.Name,
			&project.Language,
			&project.Containerization,
			&project.SrcDir,
			&project.StubDir,
		)
		projectsByIds[project.Id] = project
	}

	allTasks, _ := getTasksByProjects(ctx, utils.MapKeys(projectsByIds))
	for _, t := range allTasks {
		project := projectsByIds[t.ProjectId]
		project.Tasks = append(project.Tasks, t)
		projectsByIds[t.ProjectId] = project
	}

	return utils.MapValues(projectsByIds), nil
}

func CreateProject(ctx context.Context, p *model.Project) error {
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
	err = tx.QueryRow(ctx, "SELECT id FROM project WHERE dir = $1", p.Dir).Scan(&projectId)
	if err != nil {
		return err
	}

	moduleEntries := make([][]any, len(p.Modules))
	for i, mod := range p.Modules {
		moduleEntries[i] = []any{projectId, mod}
	}
	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"project_module"},
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

	injectFilesByTaskName := make(map[string]model.Task)
	for _, t := range p.Tasks {
		injectFilesByTaskName[t.Name] = t
	}

	injectFileEntries := [][]any{}
	rows, err := tx.Query(ctx, "SELECT id, name FROM task WHERE project_id = $1", projectId)
	for rows.Next() {
		var taskId int
		var taskName string
		err = rows.Scan(&taskId, &taskName)
		if err != nil {
			return err
		}

		task := injectFilesByTaskName[taskName]
		for name, path := range task.InjectFiles {
			injectFileEntries = append(injectFileEntries, []any{taskId, name, path})
		}
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"injectfile"},
		[]string{"task_id", "name", "path"},
		pgx.CopyFromRows(injectFileEntries),
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