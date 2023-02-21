package message

import (
	"github.com/bwmarrin/discordgo"
)

type Direct struct {
	UserID  string
	Content string
	Session *discordgo.Session
}

func (message *Direct) Send() (*discordgo.Message, error) {
	session := message.Session

	channel, err := session.UserChannelCreate(message.UserID)
	if err != nil {
		return nil, err
	}

	channelMessage := &Channel{
		ChannelID: channel.ID,
		Content:   message.Content,
		Session:   session,
	}

	return channelMessage.Send()
}
