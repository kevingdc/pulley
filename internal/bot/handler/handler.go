package handler

import "github.com/bwmarrin/discordgo"

type Handler struct{}

func Setup(session *discordgo.Session) {
	handler := &Handler{}
	session.AddHandler(handler.InteractionCreate)
}
