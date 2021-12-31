package config

import (
	"github.com/Trileon12/a/internal/app"
	"github.com/Trileon12/a/internal/storage"
	"time"
)

type Config struct {
	Storage storage.Config
	App     app.Config
}

func New() *Config {
	return &Config{
		Storage: storage.Config{MaxLength: 6},
		App: app.Config{
			HostShortURLs:   "http://localhost:8080/",
			Port:            8080,
			ShutdownTimeout: 5 * time.Second,
		},
	}
}
