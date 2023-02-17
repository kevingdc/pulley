package command

import "github.com/bwmarrin/discordgo"

type CommandHandler interface {
	HandleCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate)
}

func Handle(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	commandHandler := Resolve(interaction.ApplicationCommandData().Name)

	if commandHandler == nil {
		return
	}
	commandHandler.HandleCommand(session, interaction)
}

func Resolve(commandName string) CommandHandler {
	switch commandName {
	case SETUP_COMMAND:
		return &SetupCommand{}
	case CONNECT_COMMAND:
		return &ConnectCommand{}
	default:
		return nil
	}
}
