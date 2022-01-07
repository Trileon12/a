package app_test

import (
	"github.com/Trileon12/a/internal/app"
	"github.com/Trileon12/a/internal/config"
	"github.com/Trileon12/a/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
