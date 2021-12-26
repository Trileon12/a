package main

import (
	"github.com/Trileon12/a/internal/app"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	r.Post("/", app.GetShortURL)
	r.Get("/{ID}", app.GetFullURLByFullURL)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		println("Fatal error ", err)
	}

}
