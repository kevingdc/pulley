package interaction

import "github.com/bwmarrin/discordgo"

type Responder struct {
	Session     *discordgo.Session
	Interaction *discordgo.Interaction
}

func (res *Responder) SendInvalidChannelResponse() error {
	return res.respond("This command can only be executed in a server channel.")
}

func (res *Responder) SendMessageFailedResponse() error {
	return res.respond("Failed to send you a DM. Did you disable DM in your privacy settings?")
}

func (res *Responder) SendSentInstructionsResponse() error {
	return res.respond("I slid into your DMs. You can check the instructions there.")
}

func (res *Responder) respond(content string) error {
	return res.Session.InteractionRespond(res.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
