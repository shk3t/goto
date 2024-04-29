package router

import (
	"goto/src/config"
	"goto/src/handler"

	jwtware "github.com/gofiber/contrib/jwt"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/register", handler.Register)
	// api.Post("/login", handler.Login)

	api.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.SecretKey)},
	}))

	api.Post("/project", handler.LoadProject)
	api.Delete("/project/:id", handler.DeleteProject)

	// api.Get("/tasks", handler.GetTask)
	// api.Get("/tasks/:id", handler.GetTasks)
	// api.Post("/solution", handler.SubmitSolution)
	// api.Get("/solution/:id", handler.GetSolution)
}