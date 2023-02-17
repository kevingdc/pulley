package message

import "github.com/bwmarrin/discordgo"

type Message interface {
	Send() (*discordgo.Message, error)
}
