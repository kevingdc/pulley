package command

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

const (
	CmdSetup   = "setup"
	CmdConnect = "connect"
)

var (
	defaultSetupPermission int64 = discordgo.PermissionManageServer

	commands = []*discordgo.ApplicationCommand{
		{
			Name:                     CmdSetup,
			Description:              "Connect your Discord server to a GitHub repository",
			DefaultMemberPermissions: &defaultSetupPermission,
		},
		{
			Name:        CmdConnect,
			Description: "Connect your GitHub account to your Discord account",
		},
	}

	registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))
)

func RegisterAll(session *discordgo.Session) {
	for i, command := range commands {
		createdCommand, err := session.ApplicationCommandCreate(session.State.User.ID, "", command)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", command.Name, err)
		}
		registeredCommands[i] = createdCommand
	}
}

func DeleteAll(session *discordgo.Session) {
	for _, command := range registeredCommands {
		err := session.ApplicationCommandDelete(session.State.User.ID, "", command.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", command.Name, err)
		}
	}
}
