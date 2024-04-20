package router

import (
	"goto/src/handler"

	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api.Post("/project", handler.LoadProject)
    api.Delete("/project/:id", handler.DeleteProject)

	// api.Get("/tasks", handler.GetTask)
	// api.Get("/tasks/:id", handler.GetTasks)
	// api.Post("/solution", handler.SubmitSolution)
	// api.Get("/solution/:id", handler.GetSolution)
}