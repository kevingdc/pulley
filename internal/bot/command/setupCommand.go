package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kevingdc/pulley/internal/bot/interaction"
	"github.com/kevingdc/pulley/internal/bot/message"
)

type SetupCommand struct{}

func (command *SetupCommand) HandleCommand(session *discordgo.Session, i *discordgo.InteractionCreate) {
	responder := interaction.Responder{Session: session, Interaction: i.Interaction}

	if i.Member == nil {
		responder.SendInvalidChannelResponse()
		return
	}

	messageToSend := &message.DirectMessage{
		UserID:  i.Member.User.ID,
		Content: "Hey there! Congratulations, you just executed your first slash command",
		Session: session,
	}

	_, err := messageToSend.Send()
	if err != nil {
		responder.SendMessageFailedResponse()
		return
	}

	responder.SendSentInstructionsResponse()
}
