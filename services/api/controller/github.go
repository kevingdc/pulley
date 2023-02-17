package controller

import "github.com/gofiber/fiber/v2"

func HandleGitHubWebook(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}
