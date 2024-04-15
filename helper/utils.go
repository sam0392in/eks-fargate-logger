package helper

import (
	"os"

	"github.com/joho/godotenv"
)

var logger = Logger()

func ReadConfig(key string) string {
	err := godotenv.Load("config/.env")
	if err != nil {
		logger.Error("Error loading environment variables:", err)
		return ""
	}
	val := os.Getenv(key)
	return val
}
