package main

import (
	"context"
	"goto/src/config"
	"goto/src/database"
	"goto/src/router"

	"github.com/bytedance/sonic"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
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
		// Prefork:     true,
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	app.Use(logger.New())
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.SecretKey)},
	}))

	router.SetupRoutes(app)

	app.Listen(":3228")
}