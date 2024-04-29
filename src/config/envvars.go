package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvs() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Print("Error loading .env file")
	}
	SecretKey = os.Getenv("SECRET_KEY")
}