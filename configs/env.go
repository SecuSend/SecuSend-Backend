package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func EnvMongoURI() string {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading, using environment variables from host")
	}

	return os.Getenv("MONGOURI")
}
