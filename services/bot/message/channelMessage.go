package message

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kevingdc/pulley/pkg/app"
)

type Channel struct {
	ChannelID string
	Content   *app.MessageContent
	Session   *discordgo.Session
}

func (m *Channel) Send() (*discordgo.Message, error) {
	if m.isSimple() {
		return m.Session.ChannelMessageSend(m.ChannelID, m.Content.Body)
	}

	return m.Session.ChannelMessageSendComplex(m.ChannelID, m.toDiscordMessage())
}

func (m *Channel) isSimple() bool {
	body := m.Content.Body
	simpleMessage := app.MessageContent{
		Body: body,
	}

	return simpleMessage == *m.Content
}

func (m *Channel) toDiscordMessage() *discordgo.MessageSend {
	content := m.Content

	var header string
	if content.Header != "" {
		header = fmt.Sprintf("**%s**", content.Header)
	}

	var description string
	if content.Subtitle != "" {
		description = fmt.Sprintf("> %s", content.Subtitle)
	}

	return &discordgo.MessageSend{
		Content: header,
		Embeds: []*discordgo.MessageEmbed{{
			Type:        discordgo.EmbedTypeRich,
			Title:       content.Title,
			URL:         content.URL,
			Description: description,
			Fields: []*discordgo.MessageEmbedField{{
				Value: content.Body,
			}},
			Color: int(content.Color),
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: content.Thumbnail,
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: content.Footer,
			},
			Author: &discordgo.MessageEmbedAuthor{
				URL:     content.Author.URL,
				Name:    content.Author.Name,
				IconURL: content.Author.AvatarURL,
			},
		}},
	}
}
