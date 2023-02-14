package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kevingdc/pulley/internal/bot/config"
)

func (handler *Handler) MessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	fmt.Println(message.Content)

	if message.Content != config.BotPrefix+"ping" {
		return
	}

	// We create the private channel with the user who sent the message.
	channel, err := session.UserChannelCreate(message.Author.ID)
	if err != nil {
		fmt.Println("error creating channel:", err)
		session.ChannelMessageSend(
			message.ChannelID,
			"Something went wrong while sending the DM!",
		)
		return
	}

	// Then we send the message through the channel we created.
	_, err = session.ChannelMessageSend(channel.ID, "pong!")
	if err != nil {
		fmt.Println("error sending DM message:", err)
		session.ChannelMessageSend(
			message.ChannelID,
			"Failed to send you a DM. "+
				"Did you disable DM in your privacy settings?",
		)
	}
}
