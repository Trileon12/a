package main

import (
	"context"
	"github.com/Trileon12/a/internal/app"
	"github.com/Trileon12/a/internal/config"
	"github.com/Trileon12/a/internal/storage"
)

func main() {

	ctx := context.Background()
	conf := config.New()
	s := storage.New(&conf.Storage)
	spg := storage.NewPG(&conf.Storage)
	defer s.SaveData()
	defer spg.Close(ctx)
	application := app.New(&conf.App, s, spg)

	application.StartHTTPServer()

}
