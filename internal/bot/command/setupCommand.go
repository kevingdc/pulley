package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kevingdc/pulley/internal/bot/message"
)

type SetupCommand struct{}

func (command *SetupCommand) HandleCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	messageToSend := &message.DirectMessage{
		UserID:  interaction.Member.User.ID,
		Content: "Hey there! Congratulations, you just executed your first slash command",
		Session: session,
	}

	_, err := messageToSend.Send()

	response := "I slid into your DMs. You can check the instructions there."
	if err != nil {
		response = message.DM_ERROR_MESSAGE
	}

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
