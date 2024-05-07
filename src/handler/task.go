package handler

import (
	"context"
	"goto/src/database/query"
	"goto/src/model"
	"goto/src/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetTasks(c *fiber.Ctx) error {
	ctx := context.Background()
	user := GetCurrentUser(c)
	pager := utils.NewPager(c)

	my, _ := strconv.ParseBool(c.Query("my"))

	tasks := []model.Task{}
	if my {
		tasks = query.GetUserTasks(ctx, user.Id, pager)
	} else {
		tasks = query.GetAllTasks(ctx, pager)
	}

	return c.JSON(tasks)
}

func GetTask(c *fiber.Ctx) error {
	ctx := context.Background()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Id is not correct")
	}

	task := query.GetTaskWithStubs(ctx, id)
	if task == nil {
		return c.Status(404).SendString("Task not found")
	}

	return c.JSON(task)
}