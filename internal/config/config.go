package config

import (
	"github.com/Trileon12/a/internal/app"
	"github.com/Trileon12/a/internal/storage"
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

type Config struct {
	Storage storage.Config
	App     app.Config
}

func New() *Config {

	var cfgApp app.Config
	err := env.Parse(&cfgApp)
	if err != nil {
		log.Fatal(err)
	}
	cfgApp.ShutdownTimeout = 5 * time.Second

	var cfgStorage storage.Config
	err = env.Parse(&cfgStorage)
	if err != nil {
		log.Fatal(err)
	}
	cfgStorage.MaxLength = 6

	return &Config{
		Storage: cfgStorage,
		App:     cfgApp,
	}
}
