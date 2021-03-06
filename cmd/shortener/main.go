package main

import (
	"github.com/Trileon12/a/internal/app"
	"github.com/Trileon12/a/internal/config"
	"github.com/Trileon12/a/internal/storage"
)

func main() {

	conf := config.New()
	s := storage.MakeStorage(&conf.Storage)

	defer s.SaveData()
	defer s.Close()
	application := app.New(&conf.App, s)

	application.StartHTTPServer()

}
