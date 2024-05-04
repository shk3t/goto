package handler

import (
	"context"
	"goto/src/database/query"
	"goto/src/model"
	"goto/src/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetSolutions(c *fiber.Ctx) error {
	ctx := context.Background()
	user := GetCurrentUser(c)
	solutions := query.GetUserSolutions(ctx, user.Id)
	return c.JSON(solutions)
}

func GetSolution(c *fiber.Ctx) error {
	ctx := context.Background()
	user := GetCurrentUser(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Id is not correct")
	}

	solution := query.GetUserSolution(ctx, id, user.Id)
	if solution == nil {
		return c.Status(404).SendString("Solution not found")
	}

	return c.JSON(solution)
}

func SubmitSolution(c *fiber.Ctx) error {
	ctx := context.Background()
	user := GetCurrentUser(c)

	solution := model.SolutionBase{}
	if err := c.BodyParser(&solution); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Bad solution format")
	}

	task := query.GetTask(ctx, solution.TaskId)
	if task == nil {
		return c.Status(fiber.StatusBadRequest).SendString("Task not found")
	}

	taskFileNames := utils.MapKeys(task.Files)
	solutionFileNames := utils.MapKeys(solution.Files)
	missingFileNames := utils.Difference(taskFileNames, solutionFileNames)
	if len(missingFileNames) > 0 {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Missing files", "details": missingFileNames})
	}

	return c.SendString("OK")
}