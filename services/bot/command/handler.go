package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kevingdc/pulley/pkg/config"
)

type CommandHandler interface {
	Handle(session *discordgo.Session, interaction *discordgo.InteractionCreate)
}

func Handle(session *discordgo.Session, interaction *discordgo.InteractionCreate, config *config.Config) {
	commandHandler := Resolve(interaction.ApplicationCommandData().Name, config)

	if commandHandler == nil {
		return
	}
	commandHandler.Handle(session, interaction)
}

func Resolve(commandName string, config *config.Config) CommandHandler {
	switch commandName {
	case CmdSetup:
		return &Setup{}
	case CmdConnect:
		return &Connect{config: config}
	default:
		return nil
	}
}
