package bootstrap

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env not found")
	}

	if os.Getenv("JWT_SECRET") == "" || os.Getenv("REFRESH_SECRET") == "" {
		log.Fatal("JWT Secret cannot be empty")
	}
}
