package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/kevingdc/pulley/pkg/app"
	"github.com/kevingdc/pulley/pkg/messenger"
	"github.com/kevingdc/pulley/services/bot/command"
	"github.com/kevingdc/pulley/services/bot/handler"
	discordMessenger "github.com/kevingdc/pulley/services/bot/messenger"
)

type Bot struct {
	app     *app.App
	session *discordgo.Session
}

func New(app *app.App) (bot *Bot) {
	return &Bot{app: app}
}

func (bot *Bot) Start() {
	config := bot.app.Config
	session, err := discordgo.New("Bot " + config.BotToken)
	bot.session = session

	if err != nil {
		log.Fatal("Error creating Discord session,", err)
		return
	}

	handler.Setup(config, bot.session)

	bot.setIntents()

	err = bot.session.Open()
	if err != nil {
		log.Fatal("Error opening connection,", err)
		return
	}

	command.Setup(bot.session)
	bot.registerMessenger()

	log.Println("The bot is now running.")
}

func (bot *Bot) Stop() {
	command.DeleteAll(bot.session)
	bot.session.Close()
}

func (bot *Bot) registerMessenger() {
	messenger.Register(discordMessenger.New(bot.session))
}

func (bot *Bot) setIntents() {
	bot.session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages)
}
