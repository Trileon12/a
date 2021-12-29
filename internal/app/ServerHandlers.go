package app

import (
	"context"
	"github.com/Trileon12/a/internal/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

var strorage *storage.Storage
var conf *storage.Config

func GetShortURL(writer http.ResponseWriter, request *http.Request) {

	b, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	link := string(b)

	if link == "" {
		http.Error(writer, "Body is empty", http.StatusBadRequest)
		return
	}

	u, err := url.Parse(conf.HostShortURLs)
	if err != nil {
		http.Error(writer, "I made bad URL, sorry", http.StatusBadRequest)
	}
	u.Path = strorage.GetURLShort(link)

	shortLink := u.String()
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusCreated)

	_, err = writer.Write([]byte(shortLink))
	if err != nil {
		return
	}
}

func GetFullURLByShortURL(writer http.ResponseWriter, request *http.Request) {

	id := path.Base(request.URL.Path)

	URL, err := strorage.GetOriginalURL(id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	writer.Header().Set("Location", URL)
	writer.WriteHeader(http.StatusTemporaryRedirect)

}

func InitApp(s *storage.Storage, cfg *storage.Config) {
	strorage = s
	conf = cfg
}

func StartHTTPServer() {

	r := chi.NewRouter()

	r.Post("/", GetShortURL)
	r.Get("/{ID}", GetFullURLByShortURL)

	srv := &http.Server{Addr: ":8080", Handler: r}

	go func() {
		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			println("Fatal error ", err)
		}
	}()

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Waiting for SIGINT (kill -2)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		println("Shutdown fail", err)
	}
}
