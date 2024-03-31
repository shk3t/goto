package main

import (
	"goto/src/config"
	"goto/src/database"
	"goto/src/router"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

var dbConn *pgx.Conn

func main() {

	config.LoadConfig()
	dbConn = database.Connect()

	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	app.Use(middleware.Logger())

	router.SetupRoutes(app)

	app.Listen(":3228")
}