package handler

import (
	"context"
	"goto/src/database/query"
	"goto/src/utils"
	sc "strconv"

	"github.com/gofiber/fiber/v2"
)

func GetDelayedTasks(c *fiber.Ctx) error {
	ctx := context.Background()
	user := GetCurrentUser(c)
	pager := utils.NewPager(c)
	delayedTasks := query.GetUserDelayedTasks(ctx, user.Id, pager)
	return c.JSON(delayedTasks)
}

func GetDelayedTask(c *fiber.Ctx) error {
	ctx := context.Background()
	user := GetCurrentUser(c)

	id, err := sc.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Id is not correct")
	}

	delayedTask := query.GetUserDelayedTask(ctx, id, user.Id)
	if delayedTask == nil {
		return c.Status(404).SendString("Delayed task not found")
	}

	return c.JSON(delayedTask)
}