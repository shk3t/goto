package router

import (
	"goto/src/handler"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App) {
	// Middleware
	api := app.Group("/api")

	// routes
	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api.Post("/load-project", handler.LoadProject)

	// api.Get("/", handler.GetAllProducts)
	// api.Get("/:id", handler.GetSingleProduct)
	// api.Post("/", handler.CreateProduct)
	// api.Delete("/:id", handler.DeleteProduct)
}