package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kevingdc/pulley/pkg/config"
)

type CommandHandler interface {
	HandleCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate)
}

func Handle(session *discordgo.Session, interaction *discordgo.InteractionCreate, config *config.Config) {
	commandHandler := Resolve(interaction.ApplicationCommandData().Name, config)

	if commandHandler == nil {
		return
	}
	commandHandler.HandleCommand(session, interaction)
}

func Resolve(commandName string, config *config.Config) CommandHandler {
	switch commandName {
	case SETUP_COMMAND:
		return &SetupCommand{}
	case CONNECT_COMMAND:
		return &ConnectCommand{config: config}
	default:
		return nil
	}
}
