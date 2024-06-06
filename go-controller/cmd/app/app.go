package app

import (
	"controller_server_golang/config"
	"controller_server_golang/network"
	"controller_server_golang/repository"
	"controller_server_golang/service"
)

type App struct {
	cfg *config.Config

	repository *repository.Repository
	service    *service.Service
	network    *network.Server
}

func NewApp(cfg *config.Config) *App {
	a := &App{cfg: cfg}

	var err error
	if a.repository, err = repository.NewRepository(cfg); err != nil {
		panic(err)
	} else {
		a.service = service.NewService(a.repository)
		a.network = network.NewNetwork(a.service, cfg.Info.Port)
	}

	return a
}

func (a *App) Start() error {
	return a.network.Start()
}
