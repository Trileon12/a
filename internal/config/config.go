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
	b := getFlagValue("b")
	a := getFlagValue("a")
	f := getFlagValue("f")
	flag.Parse()
	if b != "default" {
		cfgApp.HostShortURLs = b
	}
	if a != "default" {
		cfgApp.ServerAdress = a
	}
	if f != "default" {
		cfgStorage.FilePath = f
	}

	return &Config{
		Storage: cfgStorage,
		App:     cfgApp,
	}
}

func getFlagValue(f string) string {
	var res string
	if flag.Lookup(f) == nil {
		res = *flag.String(f, "default", "")
	}
	res = flag.Lookup("b").Value.(flag.Getter).Get().(string)
	return res
}
