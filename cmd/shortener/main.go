package main

import (
	"github.com/Trileon12/a/internal/app"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	r.Post("/", app.GetShortUrl)
	r.Get("/{ID}", app.GetFullURLByFullUrl)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		println("Fatal error ", err)
	}

}
