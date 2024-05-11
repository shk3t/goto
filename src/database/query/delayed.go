package query

import (
	"context"
	db "goto/src/database"
	"goto/src/model"
	"goto/src/utils"
	"time"

	"github.com/jackc/pgx/v5"
)

func readDelayedTaskRow(row Scanable) *model.DelayedTask {
	delayedTask := model.DelayedTask{}
	err := row.Scan(
		&delayedTask.Id,
		&delayedTask.UserId,
		&delayedTask.TargetId,
		&delayedTask.Action,
		&delayedTask.Target,
		&delayedTask.Status,
		&delayedTask.Details,
		&delayedTask.UpdatedAt,
	)
	if err != nil {
		return nil
	}
	return &delayedTask
}

func readDelayedTaskRows(rows pgx.Rows) []model.DelayedTask {
	delayedTasks := []model.DelayedTask{}
	for rows.Next() {
		delayedTask := readDelayedTaskRow(rows)
		delayedTasks = append(delayedTasks, *delayedTask)
	}
	return delayedTasks
}

func GetDelayedTask(ctx context.Context, id int) *model.DelayedTask {
	row := db.ConnPool.QueryRow(ctx, "SELECT * FROM delayed_task WHERE id = $1", id)
	return readDelayedTaskRow(row)
}

func GetUserDelayedTask(ctx context.Context, id int, userId int) *model.DelayedTask {
	row := db.ConnPool.QueryRow(
		ctx,
		"SELECT * FROM delayed_task WHERE id = $1 AND user_id = $2",
		id, userId,
	)
	return readDelayedTaskRow(row)
}

func GetUserDelayedTasks(ctx context.Context, userId int, pager *utils.Pager) []model.DelayedTask {
	rows, _ := db.ConnPool.Query(
		ctx,
		"SELECT * FROM delayed_task WHERE user_id = $1`+pager.QuerySuffix()",
		userId,
	)
	return readDelayedTaskRows(rows)
}

func createDelayedTask(ctx context.Context, dt *model.DelayedTask) {
	db.ConnPool.QueryRow(
		ctx, `
        INSERT INTO delayed_task (user_id, action, target)
        VALUES ($1, $2, $3)
        RETURNING id`,
		dt.UserId, dt.Action, dt.Target,
	).Scan(&dt.Id)
}

func updateDelayedTask(ctx context.Context, dt *model.DelayedTask) {
	args := []any{dt.TargetId, dt.Status, dt.Details, dt.UpdatedAt, dt.Id}
	if dt.TargetId == 0 {
		args[0] = nil
	}

	db.ConnPool.Exec(
		ctx, `
        UPDATE delayed_task
        SET target_id = $1, status = $2, details = $3, updated_at = $4
        WHERE id = $5`,
		args...,
	)
}

func cleanupDelayedTasks(ctx context.Context, userId int) {
	db.ConnPool.Exec(
		ctx, `
        DELETE FROM delayed_task
        WHERE user_id = $1
            AND updated_at < NOW() - INTERVAL  '1 hour'`,
		userId,
	)
}

func SaveDelayedTask(ctx context.Context, dt *model.DelayedTask) *model.DelayedTask {
	cleanupDelayedTasks(ctx, dt.UserId)
	dt.UpdatedAt = time.Now()

	if dt.Id == 0 {
		dt.Status = "processing"
		createDelayedTask(ctx, dt)
	} else {
		updateDelayedTask(ctx, dt)
	}

	return dt
}