package config

import (
	"flag"
	"github.com/Trileon12/a/internal/app"
	"github.com/Trileon12/a/internal/storage"
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

var paramA string
var paramB string
var paramF string

func init() {
	if flag.Lookup("a") == nil {
		flag.StringVar(&paramA, "a", "default", "")
	}
	if flag.Lookup("b") == nil {
		flag.StringVar(&paramB, "b", "default", "")
	}
	if flag.Lookup("f") == nil {
		flag.StringVar(&paramF, "f", "default", "")
	}
}

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
	paramA = flag.Lookup("a").Value.(flag.Getter).Get().(string)
	paramB = flag.Lookup("b").Value.(flag.Getter).Get().(string)
	paramF = flag.Lookup("f").Value.(flag.Getter).Get().(string)
	flag.Parse()
	if paramB != "default" {
		cfgApp.HostShortURLs = paramB
	}
	if paramA != "default" {
		cfgApp.ServerAdress = paramA
	}
	if paramF != "default" {
		cfgStorage.FilePath = paramF
	}

	return &Config{
		Storage: cfgStorage,
		App:     cfgApp,
	}
}
