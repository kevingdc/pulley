package github

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/go-github/github"
	"github.com/kevingdc/pulley/pkg/app"
	internalhttp "github.com/kevingdc/pulley/pkg/http"
	"github.com/kevingdc/pulley/pkg/idconv"

	"github.com/kevingdc/pulley/services/api/github/event"
	eventhandler "github.com/kevingdc/pulley/services/api/github/eventHandler"

	internalgithub "github.com/kevingdc/pulley/pkg/github"
)

func handleGithubWebhook(app *app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		config := app.Config
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

		client, err := internalgithub.NewClientFromPayload(config, payloadBytes)
		if err != nil {
			log.Println(err)
			return fiber.ErrBadRequest
		}

		res, err := eventhandler.Handle(&event.Payload{
			App:     app,
			Payload: payload,
			Github:  client,
			Type:    event.Type(eventType),
		})
		if err != nil {
			log.Println(err)
			return fiber.ErrBadRequest
		}

		return c.JSON(res)
	}
}

func handleGithubOAuthRedirect(a *app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		config := a.Config
		userService := a.UserService

		if c.Query("state") != config.GithubOAuthState {
			return fiber.ErrForbidden
		}

		githubUser, err := internalgithub.GetUser(config, c.Query("code"))
		if err != nil {
			return fiber.ErrForbidden
		}

		var chatConfig app.ChatConfig
		json.Unmarshal([]byte(c.Query("chat_config")), &chatConfig)

		// TODO: Check if user with discord id already exists and also check if user with github id already exists

		newUser := &app.User{
			RepositoryID:   idconv.ToRepoID(githubUser.ID),
			RepositoryType: app.RepoGitHub,
			ChatID:         chatConfig.ChatID,
			ChatType:       app.Chat(chatConfig.ChatType),
		}

		if !userService.Exists(newUser) {
			userService.Create(newUser)
		}

		return c.Redirect("https://github.com/apps/pulley-app/installations/new")

	}
}
