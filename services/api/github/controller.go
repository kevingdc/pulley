package github

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/go-github/github"
	"github.com/kevingdc/pulley/pkg/config"
	"github.com/kevingdc/pulley/pkg/messenger"
	"github.com/kevingdc/pulley/pkg/user"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func handleGithubWebhook(config config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := &http.Request{}
		fasthttpadaptor.ConvertRequest(c.Context(), request, false)

		payloadBytes, err := github.ValidatePayload(request, []byte(config.GithubAppWebhookSecret))
		if err != nil {
			log.Println(err)
			return fiber.ErrForbidden
		}

		payload, err := github.ParseWebHook(github.WebHookType(request), payloadBytes)
		if err != nil {
			log.Println(err)
			return fiber.ErrBadRequest
		}
		log.Printf("Event type: %T\n", payload)

		messenger.Send(messenger.Message{
			User: user.User{
				RepositoryID:   "sample",
				RepositoryType: user.RepoGitHub,
				ChatID:         "156032324511727616",
				ChatType:       user.ChatDiscord,
			},
			Content: "Congrats gumana",
		})

		return c.SendString("Hello, World!")
	}
}
