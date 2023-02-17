package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kevingdc/pulley/pkg/config"
	"github.com/kevingdc/pulley/services/api"
	"github.com/kevingdc/pulley/services/bot"
)

func main() {
	config := config.Load()

	bot := bot.New(config)
	api := api.New(config)

	bot.Start()
	defer bot.Stop()

	go api.Start()
	defer api.Stop()

	waitForCloseSignal()
}

func waitForCloseSignal() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
