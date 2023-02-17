package command

import "github.com/bwmarrin/discordgo"

type ConnectCommand struct{}

func (command *ConnectCommand) HandleCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hey there! Congratulations, you just executed your first slash command",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
