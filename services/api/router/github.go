package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kevingdc/pulley/services/api/controller"
)

func registerGitHubRoutes(router fiber.Router) {
	api := router.Group("/github")

	api.Post("/webhook", controller.HandleGitHubWebook)
}
