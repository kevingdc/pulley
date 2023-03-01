package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kevingdc/pulley/pkg/app"
)

type API struct {
	app    *app.App
	server *fiber.App
}

func New(app *app.App) (api *API) {
	return &API{app: app}
}

func (api *API) Start() {
	app := fiber.New()
	api.server = app

	api.RegisterRoutes()

	log.Fatal(api.server.Listen(":" + api.app.Config.Port))
}

func (api *API) Stop() {
	api.server.Shutdown()
}
