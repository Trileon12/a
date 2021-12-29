package main

import (
	"github.com/Trileon12/a/internal/app"
	"github.com/Trileon12/a/internal/storage"
)

func main() {
	conf := storage.Config{
		MaxLength:     6,
		HostShortURLs: "http://localhost:8080/",
	}

	s := storage.New(&conf)

	app.InitApp(s, &conf)

	app.StartHTTPServer()

}
