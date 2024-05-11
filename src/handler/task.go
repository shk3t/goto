package handler

import (
	"context"
	q "goto/src/database/query"
	f "goto/src/filter"
	"goto/src/service"
	sc "strconv"

	"github.com/gofiber/fiber/v2"
)

func GetTasks(fctx *fiber.Ctx) error {
	ctx := context.Background()
	pager := service.NewPager(fctx)
	filter := f.NewTaskFilter(fctx)
	tasks := q.GetTasks(ctx, pager, filter)
	return fctx.JSON(tasks.Min())
}

func GetTask(fctx *fiber.Ctx) error {
	ctx := context.Background()

	id, err := sc.Atoi(fctx.Params("id"))
	if err != nil {
		return fctx.Status(fiber.StatusBadRequest).SendString("Id is not correct")
	}

	task := q.GetTask(ctx, id)
	if task == nil {
		return fctx.Status(404).SendString("Task not found")
	}

	return fctx.JSON(task.Private())
}