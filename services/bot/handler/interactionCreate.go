package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kevingdc/pulley/services/bot/command"
)

func (handler *Handler) InteractionCreate(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	command.Handle(session, interaction)
}
