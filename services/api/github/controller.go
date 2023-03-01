package github

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/go-github/github"
	"github.com/kevingdc/pulley/pkg/config"
	internalhttp "github.com/kevingdc/pulley/pkg/http"
	"github.com/kevingdc/pulley/pkg/user"
	"github.com/kevingdc/pulley/services/api/github/event"

	internalgithub "github.com/kevingdc/pulley/pkg/github"
)

func handleGithubWebhook(config *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := internalhttp.ToRequest(c)

		payloadBytes, err := github.ValidatePayload(request, []byte(config.GithubAppWebhookSecret))
		if err != nil {
			log.Println(err)
			return fiber.ErrForbidden
		}

		eventType := github.WebHookType(request)
		payload, err := github.ParseWebHook(eventType, payloadBytes)
		if err != nil {
			log.Println(err)
			return fiber.ErrBadRequest
		}
		log.Printf("Event type: %T\n", payload)

		client, err := internalgithub.NewClientFromPayload(config, payloadBytes)
		if err != nil {
			log.Println(err)
			return fiber.ErrBadRequest
		}

		res, err := event.Handle(&event.EventPayload{
			Config:  config,
			Payload: payload,
			Github:  client,
			Type:    event.EventType(eventType),
		})
		if err != nil {
			log.Println(err)
			return fiber.ErrBadRequest
		}

		return c.JSON(res)
	}
}

func handleGithubOAuthRedirect(config *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Query("state") != config.GithubOAuthState {
			return fiber.ErrForbidden
		}

		githubUser, err := internalgithub.GetUser(config, c.Query("code"))
		if err != nil {
			return fiber.ErrForbidden
		}

		var chatConfig user.ChatConfig
		json.Unmarshal([]byte(c.Query("chat_config")), &chatConfig)

		// TODO: Check if user with discord id already exists and also check if user with github id already exists

		newUser := &user.User{
			RepositoryID:   user.ToRepoID(githubUser.ID),
			RepositoryType: user.RepoGitHub,
			ChatID:         chatConfig.ChatID,
			ChatType:       user.Chat(chatConfig.ChatType),
		}

		if !newUser.Exists() {
			newUser.Create()
		}

		return c.Redirect("https://github.com/apps/pulley-app/installations/new")

	}
}
