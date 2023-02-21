package messenger

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kevingdc/pulley/pkg/messenger"
	"github.com/kevingdc/pulley/pkg/user"
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

func (d *DiscordMessenger) CanSend(m messenger.Message) bool {
	return m.User.ChatType == user.ChatDiscord
}

func (d *DiscordMessenger) Send(m messenger.Message) {
	messageToSend := &message.Direct{
		UserID:  m.User.ChatID,
		Content: m.Content,
		Session: d.session,
	}

	messageToSend.Send()
}
