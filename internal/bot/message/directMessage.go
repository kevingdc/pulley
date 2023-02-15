package message

import (
	"github.com/bwmarrin/discordgo"
)

type DirectMessage struct {
	UserID  string
	Content string
	Session *discordgo.Session
}

const DM_ERROR_MESSAGE = "Failed to send you a DM. Did you disable DM in your privacy settings?"

func (message *DirectMessage) Send() (*discordgo.Message, error) {
	session := message.Session

	channel, err := session.UserChannelCreate(message.UserID)
	if err != nil {
		return nil, err
	}

	channelMessage := &ChannelMessage{
		ChannelID: channel.ID,
		Content:   message.Content,
		Session:   session,
	}

	return channelMessage.Send()
}
