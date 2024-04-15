package main

import (
	"context"
	"goto/src/config"
	"goto/src/database"
	"goto/src/router"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

func main() {
	config.LoadConfig()
    config.InitDirs()

	ctx := context.Background()
	dbPool = database.Connect(ctx)

	app := fiber.New(fiber.Config{
		// Prefork:     true,
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	app.Use(logger.New())

	router.SetupRoutes(app)

	app.Listen(":3228")
}