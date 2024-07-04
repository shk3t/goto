package main

import (
	"context"
	"goto/src/config"
	"goto/src/database"
	"goto/src/router"
	"os"

	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

func main() {
	config.LoadEnvs()
	config.InitDirs()

	ctx := context.Background()
	database.Connect(ctx)
	defer database.ConnPool.Close()

	app := fiber.New(fiber.Config{
		Prefork:     false,
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	app.Use(logger.New())
	app.Use(swagger.New(swagger.Config{
		BasePath: os.Getenv("URL_PREFIX") + "/api/",
		FilePath: "./docs/swagger.json",
		Path:     "docs",
	}))

	router.SetupRoutes(app)

    app.Listen(":" + os.Getenv("PORT"))
}