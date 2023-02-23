package command

import (
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kevingdc/pulley/pkg/config"
	"github.com/kevingdc/pulley/pkg/user"
	"github.com/kevingdc/pulley/services/bot/interaction"
	"github.com/kevingdc/pulley/services/bot/message"
)

type ConnectCommand struct {
	config *config.Config
}

func (c *ConnectCommand) HandleCommand(session *discordgo.Session, i *discordgo.InteractionCreate) {
	responder := interaction.Responder{Session: session, Interaction: i.Interaction}

	if wasSentAsDM := i.Member == nil; wasSentAsDM {
		responder.SendInvalidChannelResponse()
		return
	}

	userID := i.Member.User.ID

	messageToSend := &message.Direct{
		UserID:  i.Member.User.ID,
		Content: "Hey there! To connect to your GitHub account, click on the link below, press Authorize, and select the organization you want to connect to.\nLink: " + c.generateConnectLink(userID),
		Session: session,
	}

	_, err := messageToSend.Send()
	if err != nil {
		responder.SendMessageFailedResponse()
		return
	}

	responder.SendSentInstructionsResponse()
}

func (c *ConnectCommand) generateConnectLink(userID string) string {
	clientId := c.config.GithubClientID
	state := c.config.GithubOAuthState

	chatConfig, _ := json.Marshal(user.ChatConfig{
		ChatID:   userID,
		ChatType: string(user.ChatDiscord),
	})

	redirectUri := fmt.Sprintf("%s/api/github/oauth-redirect?chat_config=%s", c.config.BaseURL, chatConfig)
	link := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&state=%s&redirect_uri=%s", clientId, state, redirectUri)
	return link
}
