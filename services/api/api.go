package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kevingdc/pulley/pkg/config"
	"github.com/kevingdc/pulley/services/api/router"
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

	router.RegisterRoutes(api.app)

	log.Fatal(api.app.Listen(":" + api.config.Port))
}

func (api *API) Stop() {
	api.app.Shutdown()
}
