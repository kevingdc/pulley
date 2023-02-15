package message

import (
	"github.com/bwmarrin/discordgo"
)

type ChannelMessage struct {
	ChannelID string
	Content   string
	Session   *discordgo.Session
}

func (message *ChannelMessage) Send() (*discordgo.Message, error) {
	session := message.Session
	return session.ChannelMessageSend(message.ChannelID, message.Content)
}
