package handler

import (
	"context"
	q "goto/src/database/query"
	"goto/src/service"
	sc "strconv"

	"github.com/gofiber/fiber/v2"
)

func GetDelayedTasks(fctx *fiber.Ctx) error {
	ctx := context.Background()
	user := service.GetCurrentUser(fctx)
	pager := service.NewPager(fctx)
	delayedTasks := q.GetUserDelayedTasks(ctx, user.Id, pager)
	return fctx.JSON(delayedTasks)
}

func GetDelayedTask(fctx *fiber.Ctx) error {
	ctx := context.Background()
	user := service.GetCurrentUser(fctx)

	id, err := sc.Atoi(fctx.Params("id"))
	if err != nil {
		return fctx.Status(fiber.StatusBadRequest).SendString("Id is not correct")
	}

	delayedTask := q.GetUserDelayedTask(ctx, id, user.Id)
	if delayedTask == nil {
		return fctx.Status(404).SendString("Delayed task not found")
	}

	return fctx.JSON(delayedTask)
}