package github

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kevingdc/pulley/pkg/app"
)

func RegisterRoutes(app *app.App, router fiber.Router) {
	apiGroup := router.Group("/github")

	apiGroup.Post("/webhook", handleGithubWebhook(app))
	apiGroup.Get("/oauth-redirect", handleGithubOAuthRedirect(app))
}
