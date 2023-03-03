package messenger

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kevingdc/pulley/pkg/app"

	"github.com/kevingdc/pulley/services/bot/message"
)

type DiscordMessenger struct {
	session *discordgo.Session
}

func New(session *discordgo.Session) *DiscordMessenger {
	return &DiscordMessenger{
		session: session,
	}
}

func (d *DiscordMessenger) CanSend(m *app.Message) bool {
	return m.User.ChatType == app.ChatDiscord
}

func (d *DiscordMessenger) Send(m *app.Message) error {
	messageToSend := &message.Direct{
		UserID:  m.User.ChatID,
		Content: m.Content,
		Session: d.session,
	}

	_, err := messageToSend.Send()

	return err
}
