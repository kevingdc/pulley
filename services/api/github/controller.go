package github

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/go-github/github"
	"github.com/kevingdc/pulley/pkg/config"
	"github.com/kevingdc/pulley/pkg/messenger"
	"github.com/kevingdc/pulley/pkg/user"
	"github.com/valyala/fasthttp/fasthttpadaptor"

	internalgithub "github.com/kevingdc/pulley/pkg/github"
)

func handleGithubWebhook(config *config.Config) fiber.Handler {
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

		newUser := &user.User{
			RepositoryID:   strconv.FormatInt(githubUser.ID, 10),
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
