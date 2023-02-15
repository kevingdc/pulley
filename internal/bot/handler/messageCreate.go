package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kevingdc/pulley/internal/bot/message"
)

func (handler *Handler) MessageCreate(session *discordgo.Session, receivedMessage *discordgo.MessageCreate) {
	if receivedMessage.Author.ID == session.State.User.ID || receivedMessage.Content != "ping" {
		return
	}

	reply := &message.DirectMessage{
		UserID:  receivedMessage.Author.ID,
		Content: "pong",
		Session: session,
	}

	_, err := reply.Send()
	if err != nil {
		session.ChannelMessageSend(
			receivedMessage.ChannelID,
			message.DM_ERROR_MESSAGE,
		)
	}
}
