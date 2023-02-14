package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	BotToken  string
	BotPrefix string
)

type Config struct {
	BotToken  string
	BotPrefix string
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	BotToken = os.Getenv("BOT_TOKEN")
	BotPrefix = os.Getenv("BOT_PREFIX")
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("BOT_TOKEN")
	botPrefix := os.Getenv("BOT_PREFIX")

	config := Config{BotToken: botToken, BotPrefix: botPrefix}

	return config
}
