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
	api.Post("/login", handler.Login)

	api.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.SecretKey)},
	}))

	api.Post("/projects", handler.LoadProject)
	api.Delete("/projects/:id", handler.DeleteProject)

	api.Get("/projects", handler.GetProjects)
    api.Get("/projects/:id", handler.GetProject)

	api.Get("/tasks", handler.GetTasks)
    api.Get("/tasks/:id", handler.GetTask)

	api.Get("/solutions", handler.GetSolutions)
    api.Get("/solutions/:id", handler.GetSolution)
}