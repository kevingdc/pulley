package bot

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/kevingdc/pulley/pkg/config"
	"github.com/kevingdc/pulley/services/bot/command"
	"github.com/kevingdc/pulley/services/bot/handler"
)

type Bot struct {
	config  config.Config
	session *discordgo.Session
}

func New(config config.Config) (bot *Bot) {
	return &Bot{config: config}
}

func (bot *Bot) Start() {
	session, err := discordgo.New("Bot " + bot.config.BotToken)
	bot.session = session

	if err != nil {
		log.Fatal("error creating Discord session,", err)
		return
	}

	handler.Setup(bot.session)

	setIntents(bot.session)

	err = bot.session.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
		return
	}

	command.Setup(session)

	fmt.Println("The bot is now running.")
}

func (bot *Bot) Stop() {
	command.DeleteAll(bot.session)
	bot.session.Close()
}

func setIntents(session *discordgo.Session) {
	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages)
}
