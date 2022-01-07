package config

import (
	"flag"
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

	cfgApp, err := GetAppConfig()
	cfgStorage := GetStorageConfig(err)

	return &Config{
		Storage: cfgStorage,
		App:     cfgApp,
	}
}

func GetStorageConfig(err error) storage.Config {
	cfgStorage := storage.Config{}
	err = env.Parse(&cfgStorage)
	if err != nil {
		log.Fatal(err)
	}

	if isFlagPassed("f") {
		flag.StringVar(&cfgStorage.FilePath, "f", "", "path to crazy db")
	}

	cfgStorage.MaxLength = 6
	return cfgStorage
}

func GetAppConfig() (app.Config, error) {
	cfgApp := app.Config{}

	err := env.Parse(&cfgApp)
	if err != nil {
		log.Fatal(err)
	}
	if isFlagPassed("b") {
		flag.StringVar(&cfgApp.HostShortURLs, "b", "http://localhost:8080/", "base url")
	}
	if isFlagPassed("a") {
		flag.StringVar(&cfgApp.ServerAdress, "a", ":8080", "server url")
	}

	cfgApp.ShutdownTimeout = 5 * time.Second
	return cfgApp, err
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
