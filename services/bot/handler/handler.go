package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kevingdc/pulley/pkg/config"
)

type Handler struct {
	config *config.Config
}

func Setup(config *config.Config, session *discordgo.Session) {
	handler := &Handler{config: config}
	session.AddHandler(handler.InteractionCreate)
}
