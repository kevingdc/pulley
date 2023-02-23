package github

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kevingdc/pulley/pkg/config"
)

func RegisterRoutes(config *config.Config, router fiber.Router) {
	apiGroup := router.Group("/github")

	apiGroup.Post("/webhook", handleGithubWebhook(config))
	apiGroup.Get("/oauth-redirect", handleGithubOAuthRedirect(config))
}
