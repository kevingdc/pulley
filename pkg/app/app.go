package app

import "github.com/kevingdc/pulley/pkg/config"

type App struct {
	Config      *config.Config
	UserService UserService
}
