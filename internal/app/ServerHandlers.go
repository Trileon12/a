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
	"strconv"
	"syscall"
	"time"
)

type Config struct {
	HostShortURLs   string
	Port            int
	ShutdownTimeout time.Duration
}

type App struct {
	conf    *Config
	storage *storage.Storage
}

func New(conf *Config, storage *storage.Storage) *App {
	return &App{conf, storage}
}

// Get short URL for full url
func (a *App) GetShortURL(writer http.ResponseWriter, request *http.Request) {

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

	u, err := url.Parse(a.conf.HostShortURLs)
	if err != nil {
		http.Error(writer, "I made bad URL, sorry", http.StatusBadRequest)
	}
	u.Path = a.storage.GetURLShort(link)

	shortLink := u.String()
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusCreated)

	_, err = writer.Write([]byte(shortLink))
	if err != nil {
		http.Error(writer, "I have short URL, but not for you", http.StatusBadRequest)
		return
	}
}

// Get full URL by short URL
func (a *App) GetFullURLByShortURL(writer http.ResponseWriter, request *http.Request) {

	id := path.Base(request.URL.Path)

	URL, err := a.storage.GetOriginalURL(id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	writer.Header().Set("Location", URL)
	writer.WriteHeader(http.StatusTemporaryRedirect)

}

func (a *App) StartHTTPServer() {

	srv := a.Routing()

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

	ctx, cancel := context.WithTimeout(context.Background(), a.conf.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		println("Shutdown fail", err)
	}
}

func (a *App) Routing() *http.Server {
	r := chi.NewRouter()

	r.Post("/", a.GetShortURL)
	r.Get("/{ID}", a.GetFullURLByShortURL)

	srv := &http.Server{Addr: ":" + strconv.Itoa(a.conf.Port), Handler: r}
	return srv
}
