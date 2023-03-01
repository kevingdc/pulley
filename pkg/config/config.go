package config

import (
	"encoding/base64"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL string

	BotToken string

	Port                    string
	BaseURL                 string
	GithubBaseURL           string
	GithubAppID             int64
	GithubAppWebhookSecret  string
	GithubClientID          string
	GithubClientSecret      string
	GithubAppPrivateKeyPath string
	GithubOAuthState        string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	privateKeyPath, err := createTempFileForPrivateKey(os.Getenv("GITHUB_APP_PRIVATE_KEY"))
	if err != nil {
		log.Fatal("Error creating temp file for GitHub app private key")
	}

	githubAppID, err := strconv.ParseInt(os.Getenv("GITHUB_APP_ID"), 10, 64)
	if err != nil {
		log.Fatal("Error parsing GitHub App ID")
	}

	config := &Config{
		DBURL: os.Getenv("DB_URL"),

		BotToken: os.Getenv("BOT_TOKEN"),

		Port:                    os.Getenv("PORT"),
		BaseURL:                 os.Getenv("BASE_URL"),
		GithubBaseURL:           os.Getenv("GITHUB_BASE_URL"),
		GithubAppID:             githubAppID,
		GithubAppWebhookSecret:  os.Getenv("GITHUB_APP_WEBHOOK_SECRET"),
		GithubClientID:          os.Getenv("GITHUB_CLIENT_ID"),
		GithubClientSecret:      os.Getenv("GITHUB_CLIENT_SECRET"),
		GithubAppPrivateKeyPath: privateKeyPath,
		GithubOAuthState:        os.Getenv("GITHUB_OAUTH_STATE"),
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
