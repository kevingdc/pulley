package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kevingdc/pulley/pkg/app"
	"github.com/kevingdc/pulley/pkg/config"
	"github.com/kevingdc/pulley/pkg/db/postgres"
	"github.com/kevingdc/pulley/services/api"
	"github.com/kevingdc/pulley/services/bot"
)

func main() {
	config := config.Load()

	postgres.Connect(config.DBURL)
	defer postgres.Close()

	newApp := &app.App{
		Config:      config,
		UserService: &postgres.UserService{DB: postgres.DB},
	}

	bot := bot.New(newApp)
	bot.Start()
	defer bot.Stop()

	api := api.New(newApp)
	go api.Start()
	defer api.Stop()

	waitForCloseSignal()
}

func waitForCloseSignal() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
