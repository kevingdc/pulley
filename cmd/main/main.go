package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kevingdc/pulley/pkg/config"
	"github.com/kevingdc/pulley/pkg/db"
	"github.com/kevingdc/pulley/services/api"
	"github.com/kevingdc/pulley/services/bot"
)

func main() {
	config := config.Load()

	db.Connect(config)
	defer db.Close()

	bot := bot.New(config)
	bot.Start()
	defer bot.Stop()

	api := api.New(config)
	go api.Start()
	defer api.Stop()

	waitForCloseSignal()
}

func waitForCloseSignal() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
