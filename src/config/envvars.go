package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadEnvs() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Print("Error loading .env file")
	}
}