package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Print("Error loading .env file")
	}
}