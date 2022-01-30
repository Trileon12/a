package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Trileon12/a/internal/Middleware"
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

type Config struct {
	HostShortURLs   string        `env:"BASE_URL" envDefault:"http://localhost:8080/"`
	ServerAddress   string        `env:"SERVER_ADDRESS" envDefault:":8080"`
	ShutdownTimeout time.Duration `env:"ShutdownTimeout" envDefault:"5s"`
}

type App struct {
	conf    *Config
	storage *storage.Storage
}

type ShortURLRequest struct {
	URL string `json:"url"`
}

type ShortURLResponse struct {
	Result string `json:"result"`
}

func New(conf *Config, storage *storage.Storage) *App {
	return &App{conf, storage}
}

// GetShortURL Get short URL for full url
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

	userID := request.Header.Get("userID")

	u.Path = a.storage.GetURLShort(link, userID)

	shortLink := u.String()
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusCreated)

	_, err = writer.Write([]byte(shortLink))
	if err != nil {
		http.Error(writer, "I have short URL, but not for you", http.StatusBadRequest)
		return
	}
}

// GetShortURLJson Get short URL for full url JSON format
func (a *App) GetShortURLJson(writer http.ResponseWriter, request *http.Request) {

	var b ShortURLRequest

	if err := json.NewDecoder(request.Body).Decode(&b); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if b.URL == "" {
		http.Error(writer, "URL is empty", http.StatusBadRequest)
		return
	}

	u, err := url.Parse(a.conf.HostShortURLs)
	if err != nil {
		http.Error(writer, "I made bad URL, sorry", http.StatusBadRequest)
	}
	userID := request.Header.Get("userID")
	u.Path = a.storage.GetURLShort(b.URL, userID)

	resp := ShortURLResponse{}
	resp.Result = u.String()

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(writer).Encode(resp)
	if err != nil {
		http.Error(writer, "I have short URL, but not for you", http.StatusBadRequest)
		return
	}
}

func (a *App) GetUserURLs(writer http.ResponseWriter, request *http.Request) {
	userID := request.Header.Get("userID")
	userURLS := a.storage.GetUserURLS(userID)

	if len(userURLS) == 0 {
		writer.WriteHeader(http.StatusNoContent)
	} else {
		for i, _ := range userURLS {
			userURLS[i].ShortURL = a.conf.HostShortURLs + userURLS[i].ShortURL
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(userURLS)
}

// GetFullURLByShortURL Get full URL by short URL
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
			fmt.Println("Fatal  error ", err)
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

	r.Use(middleware.UnzipHandle)
	r.Use(middleware.ZipHandle)
	r.Use(middleware.SetUserIDCookieHandle)
	r.Post("/", a.GetShortURL)
	r.Post("/api/shorten", a.GetShortURLJson)
	r.Get("/{ID}", a.GetFullURLByShortURL)
	r.Get("/user/urls", a.GetUserURLs)

	srv := &http.Server{Addr: a.conf.ServerAddress, Handler: r}
	return srv
}
