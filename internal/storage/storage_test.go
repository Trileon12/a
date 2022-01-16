package storage_test

import (
	"github.com/Trileon12/a/internal/config"
	"github.com/Trileon12/a/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var conf *config.Config

func TestMain(m *testing.M) {
	conf = config.New()
	appTsts := m.Run()

	os.Exit(appTsts)

}

func TestGetOriginalURL(t *testing.T) {
	type args struct {
		shortURL string
	}

	s := storage.New(&conf.Storage)

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "DoesntURL",
			args: args{shortURL: "www.google.com"},
			want: "www.google.com",
		},
		{
			name: "Some text ",
			args: args{shortURL: "stylish exterior belies a growing darkness that has long been smoldering within, revealing that true monsters are made, not born"},
			want: "www.google.con",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			shortURL := s.GetURLShort(tt.args.shortURL)

			got, err := s.GetOriginalURL(shortURL)

			require.NoError(t, err)
			require.Equal(t, tt.args.shortURL, got)

		})
	}
}

func TestGetOriginalURLErr(t *testing.T) {
	type args struct {
		shortURL string
	}

	s := storage.New(&conf.Storage)

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Std str",
			args: args{shortURL: storage.RandString(5)},
		},
		{
			name: "big str",
			args: args{shortURL: storage.RandString(50)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := s.GetOriginalURL(tt.args.shortURL)

			require.Error(t, err)

		})
	}
}

func TestGetURLShort(t *testing.T) {
	type args struct {
		originalURL string
	}

	s := storage.New(&conf.Storage)

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "std url",
			args: args{originalURL: "www.google.com"},
			want: "[[:alpha:]]{6}$",
		},
		{
			name: "bad url",
			args: args{originalURL: "bla bla bla"},
			want: "[[:alpha:]]{6}$",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.GetURLShort(tt.args.originalURL)
			assert.Regexp(t, tt.want, got, "Random Str doesn't match the pattern")
		})
	}
}

func TestRandString(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "std string ",
			args: args{
				n: 6,
			},
			want: "[[:alpha:]]{6}$",
		},
		{
			name: "zero sting",
			args: args{
				n: 0,
			},
			want: "[[:alpha:]]{0}$",
		},

		{
			name: "big string",
			args: args{
				n: 200,
			},
			want: "[[:alpha:]]{200}$",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := storage.RandString(tt.args.n)

			assert.Regexp(t, tt.want, got, "Random Str doesn't match the pattern")
		})
	}
}
