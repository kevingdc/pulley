package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kevingdc/pulley/pkg/config"
)

type API struct {
	config config.Config
	app    *fiber.App
}

func New(config config.Config) (api *API) {
	return &API{config: config}
}

func (api *API) Start() {
	app := fiber.New()
	api.app = app

	api.RegisterRoutes()

	log.Fatal(api.app.Listen(":" + api.config.Port))
}

func (api *API) Stop() {
	api.app.Shutdown()
}
