package message

import (
	"github.com/bwmarrin/discordgo"
)

type Channel struct {
	ChannelID string
	Content   string
	Session   *discordgo.Session
}

func (message *Channel) Send() (*discordgo.Message, error) {
	session := message.Session
	return session.ChannelMessageSend(message.ChannelID, message.Content)
}
