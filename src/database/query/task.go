package query

import (
	"context"
	db "goto/src/database"
	"goto/src/model"
	"goto/src/utils"
)

func getTasksByProjects(ctx context.Context, projectIds []int) ([]model.Task, error) {
	tasksByIds := make(map[int]model.Task)

	rows, _ := db.ConnPool.Query(
		ctx,
		"SELECT * FROM task WHERE project_id = ANY ($1)",
		projectIds,
	)

	for rows.Next() {
		task := model.Task{}
		rows.Scan(
			&task.Id,
			&task.ProjectId,
			&task.Name,
			&task.Description,
			&task.RunTarget,
		)
		task.InjectFiles = make(map[string]string)
		tasksByIds[task.Id] = task
	}

	allInjectFiles, _ := getInjectFilesByTasks(ctx, utils.MapKeys(tasksByIds))
	for _, ifl := range allInjectFiles {
		tasksByIds[ifl.TaskId].InjectFiles[ifl.Name] = ifl.Path
	}

	return utils.MapValues(tasksByIds), nil
}

func getInjectFilesByTasks(ctx context.Context, taskIds []int) ([]model.InjectFile, error) {
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

	return allInjectFiles, nil
}