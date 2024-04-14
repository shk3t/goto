package main

import (
	"context"
	"goto/src/config"
	"goto/src/database"
	"goto/src/router"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/jackc/pgx/v5"
)

var dbConn *pgx.Conn

func main() {
	ctx := context.Background()

	config.LoadConfig()
	dbConn = database.Connect(ctx)

	app := fiber.New(fiber.Config{
		// Prefork:     true,
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	app.Use(logger.New())

	router.SetupRoutes(app)

	app.Listen(":3228")
}