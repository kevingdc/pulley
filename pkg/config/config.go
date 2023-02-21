package config

import (
	"encoding/base64"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string

	Port                    string
	GithubBaseURL           string
	GithubAppID             string
	GithubAppWebhookSecret  string
	GithubClientSecret      string
	GithubAppPrivateKeyPath string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	privateKeyPath, err := createTempFileForPrivateKey(os.Getenv("GITHUB_APP_PRIVATE_KEY"))
	if err != nil {
		log.Fatal("Error creating temp file for GitHub app private key")
	}

	config := Config{
		BotToken: os.Getenv("BOT_TOKEN"),

		Port:                    os.Getenv("PORT"),
		GithubBaseURL:           os.Getenv("GITHUB_BASE_URL"),
		GithubAppID:             os.Getenv("GITHUB_APP_ID"),
		GithubAppWebhookSecret:  os.Getenv("GITHUB_APP_WEBHOOK_SECRET"),
		GithubClientSecret:      os.Getenv("GITHUB_CLIENT_SECRET"),
		GithubAppPrivateKeyPath: privateKeyPath,
	}

	return config
}

func createTempFileForPrivateKey(encodedPrivateKey string) (string, error) {
	file, err := os.CreateTemp("", "github-app-private-key")

	if err != nil {
		return "", err
	}

	privateKey, err := base64.StdEncoding.DecodeString(encodedPrivateKey)
	if err != nil {
		return "", err
	}

	_, err = file.WriteString(string(privateKey))
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}
