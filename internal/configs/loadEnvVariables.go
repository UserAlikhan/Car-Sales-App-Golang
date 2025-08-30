package configs

import (
	"log"

	"github.com/joho/godotenv"
)

// function that checks if we
// have an access to the .env file
func LoadEnvVariables() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
