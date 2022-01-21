package app

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//Если gzip не поддерживатеся, то ничего не делем
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := zipping(w)
		if err != nil {
			http.Error(w, "Un zip err", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)

	})
}

func zipping(w http.ResponseWriter) (*gzip.Writer, error) {
	// создаём gzip.Writer поверх текущего w
	gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	if err != nil {
		io.WriteString(w, err.Error())
		return nil, err
	}
	defer gz.Close()

	w.Header().Set("Content-Encoding", "gzip")
	return gz, nil
}
