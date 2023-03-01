package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kevingdc/pulley/services/api/github"
)

func (api *API) RegisterRoutes() {
	apiGroup := api.server.Group("/api")

	apiGroup.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	github.RegisterRoutes(api.app, apiGroup)
}
