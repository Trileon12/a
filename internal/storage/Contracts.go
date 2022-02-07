package storage

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type Config struct {
	MaxLength int    `env:"MaxLength" envDefault:"6"`
	FilePath  string `env:"FILE_STORAGE_PATH"`
	DBAddress string `env:"DATABASE_DSN"`
}

func (c *Config) IsDBDefined() bool {
	return c.DBAddress != ""
}

type URLsType map[string]string
type URLPair struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type UserURLList struct {
	UserID string
	URLs   []string
}

type ShortURLItemRequest struct {
	OriginalURL   string `json:"original_url"`
	CorrelationID string `json:"correlation_id"`
}

type ShortURLItemResponse struct {
	ShortURL      string `json:"short_url"`
	CorrelationID string `json:"correlation_id"`
}

type Storage interface {
	GetURLShort(originalURL string, userID string) (string, error)
	GetURLsShort(originalURL []ShortURLItemRequest, userID string, host string) ([]ShortURLItemResponse, error)
	GetUserURLS(userID string) []URLPair
	DeleteURLS(userID string, URLs []string)
	Ping(ctx context.Context) error
	GetOriginalURL(shortURL string, userID string) (string, error)
	Close()
	SaveData()
}

func MakeStorage(conf *Config) Storage {
	if conf.IsDBDefined() {
		return NewStorageDB(conf)
	} else {
		return NewStorageMap(conf)
	}
}

func ExtractJSONURLData(fileName string, s *URLsType) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	_ = json.Unmarshal(byteValue, &s)
}

var ErrDuplicateOriginalURL = errors.New("дублированная ссылка")
var ErrURLDeleted = errors.New("URL был удален")
