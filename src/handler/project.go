package handler

import (
	"goto/src/model"

	"github.com/gofiber/fiber/v2"
)

func LoadProject(c *fiber.Ctx) error {
	project := new(model.Project)

	if err := c.BodyParser(project); err != nil {
		return c.Status(500).JSON(err)
	}

	return c.JSON(&fiber.Map{"status": "success"})
}