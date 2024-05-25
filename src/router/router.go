package router

import (
	"goto/src/config"
	"goto/src/handler"

	jwtware "github.com/gofiber/contrib/jwt"

	"github.com/gofiber/fiber/v2"
)

// @title Goto
// @description Web app for code challenges with any environments
// @contact.name Goto GitHub
// @contact.url http://github.com/shk3t/goto
// @host localhost:3228
// @basePath /api/
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Prepend your JWT key with `Bearer`
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/register", handler.Register)
	api.Post("/login", handler.Login)

	api.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.SecretKey)},
	}))

	api.Get("/projects", handler.GetProjects)
	api.Get("/projects/:id", handler.GetProject)

	api.Get("/tasks", handler.GetTasks)
	api.Get("/tasks/:id", handler.GetTask)

	api.Get("/solutions", handler.GetSolutions)
	api.Get("/solutions/:id", handler.GetSolution)

	api.Post("/projects", handler.LoadProject)
	api.Put("/projects/:id", handler.LoadProject)
	api.Delete("/projects/:id", handler.DeleteProject)
	api.Post("/solutions", handler.SubmitSolution)

	api.Get("/delayed-tasks", handler.GetDelayedTasks)
	api.Get("/delayed-tasks/:id", handler.GetDelayedTask)
}