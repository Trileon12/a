package config

import (
	"flag"
	"github.com/Trileon12/a/internal/app"
	"github.com/Trileon12/a/internal/storage"
	"github.com/caarlos0/env/v6"
	"log"
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

	var cfgStorage storage.Config
	err = env.Parse(&cfgStorage)
	if err != nil {
		log.Fatal(err)
	}
	cfgStorage.MaxLength = 6

	flag.StringVar(&cfgApp.ServerAddress, "a", cfgApp.ServerAddress, "port to run server")
	flag.StringVar(&cfgApp.HostShortURLs, "b", cfgApp.HostShortURLs, "base URL for shorten URL response")
	flag.StringVar(&cfgStorage.FilePath, "f", cfgStorage.FilePath, "file to store shorten URLs")
	flag.Parse()

	return &Config{
		Storage: cfgStorage,
		App:     cfgApp,
	}
}
