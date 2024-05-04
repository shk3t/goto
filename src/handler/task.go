package handler

import (
	"context"
	"goto/src/database/query"
	"goto/src/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetTasks(c *fiber.Ctx) error {
	ctx := context.Background()
	user := GetCurrentUser(c)

	my, _ := strconv.ParseBool(c.Query("my"))

	var tasks []model.Task
	if my {
		tasks = query.GetUserTasks(ctx, user.Id)
	} else {
		tasks = query.GetAllTasks(ctx)
	}

	return c.JSON(tasks)
}

func GetTask(c *fiber.Ctx) error {
	ctx := context.Background()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Id is not correct")
	}

	task := query.GetTask(ctx, id)
	if task == nil {
		return c.Status(404).SendString("Task not found")
	}

	return c.JSON(task)
}