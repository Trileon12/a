package Middleware

import (
	"compress/gzip"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

var ErrWrongSign = errors.New("wrong sign")

//unzip requset
func UnzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {

		if request.Header.Get(`Content-Encoding`) == `gzip` {
			gz, err := gzip.NewReader(request.Body)
			if err != nil {
				http.Error(response, err.Error(), http.StatusBadRequest)
				return
			}
			request.Body = gz
		}
		next.ServeHTTP(response, request)

	})
}

//zip response
func ZipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if !strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(response, request)
			return
		}
		gzipResponse, err := gzip.NewWriterLevel(response, gzip.BestSpeed)
		if err != nil {
			io.WriteString(response, err.Error())
			return
		}
		defer gzipResponse.Close()

		response.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipResponseWriter{ResponseWriter: response, Writer: gzipResponse}, request)
	})
}

//userCookie
func SetUserIDCookieHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {

		userIDCookie, err := request.Cookie("token")
		createNewIdentity := false

		if err == http.ErrNoCookie {
			createNewIdentity = true
		} else {
			id, err := checkSign(userIDCookie.Value)
			if err == ErrWrongSign {
				createNewIdentity = true
			}
			request.Header.Set("userID", strconv.Itoa(int(id)))
		}

		if createNewIdentity {
			id, err := genId()
			if err != nil {
				http.Error(response, err.Error(), http.StatusBadRequest)
				return
			}
			token, err := Sign(id)
			if err != nil {
				http.Error(response, err.Error(), http.StatusBadRequest)
				return
			}
			userIDCookie = &http.Cookie{Name: "tocken", Value: token}
			request.AddCookie(userIDCookie)
		}

		next.ServeHTTP(response, request)
		http.SetCookie(response, userIDCookie)
	})
}

func genId() ([]byte, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func checkSign(msg string) (uint32, error) {
	var (
		data []byte
		id   uint32
		sign []byte
	)

	data, _ = hex.DecodeString(msg)
	id = binary.BigEndian.Uint32(data[:4])
	h := hmac.New(sha256.New, []byte("secret key"))
	h.Write(data[:4])
	sign = h.Sum(nil)

	if hmac.Equal(sign, data[4:]) {
		return id, nil
	} else {
		return 0, ErrWrongSign
	}
}

func Sign(id []byte) (string, error) {
	h := hmac.New(sha256.New, []byte("secret key"))
	h.Write(id)
	dst := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(append(id, dst...)), nil
}
