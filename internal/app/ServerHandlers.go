package app

import (
	"encoding/json"
	"github.com/Trileon12/a/internal/storage"
	"io"
	"net/http"
	"path"
)

var HostShortURLs string

func init() {
	HostShortURLs = "http://localhost:8080/"
}

type ShortLink string

func GetShortURL(writer http.ResponseWriter, request *http.Request) {

	b, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	link := string(b)

	if link == "" {
		http.Error(writer, "Body is empty", http.StatusInternalServerError)
		return
	}
	shortLink := HostShortURLs + storage.GetURLShort(link)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	resp, err := json.Marshal(shortLink)
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}
	_, err = writer.Write(resp)
	if err != nil {
		return
	}
}

func GetFullURLByFullURL(writer http.ResponseWriter, request *http.Request) {

	id := path.Base(request.URL.Path)

	URL, err := storage.GetOriginalURL(id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	writer.Header().Set("Location", URL)
	writer.WriteHeader(http.StatusTemporaryRedirect)

}
