package console

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("DOTENV: No .env file found or error reading file.")
	} else {
		log.Printf("DOTENV: .env found.")
	}
}
