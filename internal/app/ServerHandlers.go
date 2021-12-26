package app

import (
	"github.com/Trileon12/a/internal/storage"
	"io"
	"net/http"
	"path"
)

var HostShortURLs string

func init() {
	HostShortURLs = "http://localhost:8080/"
}

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
	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusCreated)
	_, err = writer.Write([]byte(shortLink))
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
	writer.WriteHeader(http.StatusMovedPermanently)

}
