package config

import (
	"flag"
	"github.com/Trileon12/a/internal/app"
	"github.com/Trileon12/a/internal/storage"
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

var a string
var b string
var f string

func init() {
	if flag.Lookup("a") == nil {
		flag.StringVar(&a, "a", "default", "")
	}
	if flag.Lookup("b") == nil {
		flag.StringVar(&b, "b", "default", "")
	}
	if flag.Lookup("f") == nil {
		flag.StringVar(&f, "f", "default", "")
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
	a = flag.Lookup("a").Value.(flag.Getter).Get().(string)
	b = flag.Lookup("b").Value.(flag.Getter).Get().(string)
	f = flag.Lookup("f").Value.(flag.Getter).Get().(string)
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
