package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
	Port     string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := Config{
		BotToken: os.Getenv("BOT_TOKEN"),
		Port:     os.Getenv("PORT"),
	}

	return config
}
