package app_test

import (
	"bytes"
	"encoding/json"
	"github.com/Trileon12/a/internal/app"
	"github.com/Trileon12/a/internal/config"
	"github.com/Trileon12/a/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var host string = "http://localhost:8080/"

type want struct {
	contentType string
	statusCode  int
	regexpLink  string
}

type request struct {
	method      string
	url         string
	originalURL string
	jsonFormat  bool
}

type tstRequest struct {
	nameTest string
	request  request
	want1    want
}

func TestGetShortURL(t *testing.T) {

	conf := config.New()
	s := storage.New(&conf.Storage)
	application := app.New(&conf.App, s)

	tests := []tstRequest{
		{

			nameTest: "Get standard URL _",
			request: request{
				method:      http.MethodPost,
				url:         "/",
				originalURL: "www.google.com",
			},
			want1: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusCreated,
				regexpLink:  "^" + host + "[[:alpha:]]{6}$",
			},
		},
		{

			nameTest: "Get standard URL JSON _",
			request: request{
				method:      http.MethodPost,
				url:         "/api/shorten",
				originalURL: "www.google.com",
				jsonFormat:  true,
			},
			want1: want{
				contentType: "application/json",
				statusCode:  http.StatusCreated,
				regexpLink:  "^" + host + "[[:alpha:]]{6}$",
			},
		},
		{

			nameTest: "Get big URL",
			request: request{
				method:      http.MethodPost,
				url:         "/",
				originalURL: "https://www.google.com/maps/place/%D0%A2%D0%B0%D0%B4%D0%B5%D0%B1%D1%8F-%D0%AF%D1%85%D0%B0,+%D0%AF%D0%BC%D0%B0%D0%BB%D0%BE-%D0%9D%D0%B5%D0%BD%D0%B5%D1%86%D0%BA%D0%B8%D0%B9+%D0%B0%D0%B2%D1%82%D0%BE%D0%BD%D0%BE%D0%BC%D0%BD%D1%8B%D0%B9+%D0%BE%D0%BA%D1%80%D1%83%D0%B3,+629705/@70.3779226,74.1234431,15z/data=!3m1!4b1!4m5!3m4!1s0x4497ae2225174a49:0xbf4bb88041f8a6f3!8m2!3d70.3779692!4d74.132309",
			},
			want1: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusCreated,
				regexpLink:  "^" + host + "[[:alpha:]]{6}$",
			}},
		{

			nameTest: "Without URL",
			request: request{
				method:      http.MethodPost,
				url:         "/",
				originalURL: "",
			},
			want1: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				regexpLink:  "",
			}},
	}

	r := chi.NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.nameTest, func(t *testing.T) {

			var request *http.Request
			var result *http.Response
			if tt.request.jsonFormat {
				user := &app.ShortURLRequest{Url: tt.request.originalURL}
				b, _ := json.Marshal(user)
				request = httptest.NewRequest(tt.request.method, tt.request.url, bytes.NewBuffer(b))
				result = SendRequest(request, application.GetShortURLJson)
			} else {
				request = httptest.NewRequest(tt.request.method, tt.request.url, strings.NewReader(tt.request.originalURL))
				result = SendRequest(request, application.GetShortURL)
			}

			assert.Equal(t, tt.want1.statusCode, result.StatusCode)
			if result.StatusCode == http.StatusCreated {
				assert.Equal(t, tt.want1.contentType, result.Header.Get("Content-Type"))

				var link string
				if tt.request.jsonFormat {
					var b app.ShortURLResponse
					if err := json.NewDecoder(result.Body).Decode(&b); err != nil {
						assert.Error(t, err, "Can't get JSON response")
					}
					link = b.Result
				} else {
					body, err := ioutil.ReadAll(result.Body)
					require.NoError(t, err)
					defer result.Body.Close()
					link = string(body)
				}
				assert.Regexp(t, tt.want1.regexpLink, link, "Short URL doesn't match the pattern")

				requestShortURL := httptest.NewRequest(http.MethodGet, link, strings.NewReader(tt.request.originalURL))
				resultShort := SendRequest(requestShortURL, application.GetFullURLByShortURL)

				assert.Equal(t, http.StatusTemporaryRedirect, resultShort.StatusCode)
				assert.Equal(t, tt.request.originalURL, resultShort.Header.Get("Location"), "Sent and got link is different")
				defer resultShort.Body.Close()
				//if err != nil {
				//	assert.Error(t, err, "Error on read body after get short link")
				//}
			}

		})
	}
}

func TestShortURL(t *testing.T) {

	conf := config.New()
	s := storage.New(&conf.Storage)
	application := app.New(&conf.App, s)

	tests := []tstRequest{
		{

			nameTest: "Tst not found URL",
			request: request{
				method:      http.MethodGet,
				url:         host + "/qGHrty",
				originalURL: "www.google.com",
			},
			want1: want{

				statusCode: http.StatusNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nameTest, func(t *testing.T) {

			request := httptest.NewRequest(tt.request.method, tt.request.url, strings.NewReader(tt.request.originalURL))
			result := SendRequest(request, application.GetFullURLByShortURL)

			assert.Equal(t, tt.want1.statusCode, result.StatusCode)

			defer result.Body.Close()
			//if err != nil {
			//	assert.Error(t, err, "Error on read body after get short link")
			//}

		})
	}
}

func SendRequest(request *http.Request, f func(writer http.ResponseWriter, request *http.Request)) *http.Response {
	w := httptest.NewRecorder()
	h := http.HandlerFunc(f)
	h.ServeHTTP(w, request)
	result := w.Result()
	return result
}
