package bot

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/kevingdc/pulley/internal/bot/command"
	"github.com/kevingdc/pulley/internal/bot/config"
	"github.com/kevingdc/pulley/internal/bot/handler"
)

func Start() {
	session, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	handler.Setup(session)

	setIntents(session)

	err = session.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	defer session.Close()

	command.Setup(session)
	defer command.DeleteAll(session)

	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	waitForCloseSignal()
}

func setIntents(session *discordgo.Session) {
	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages)
}

func waitForCloseSignal() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
