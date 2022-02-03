package storage

import (
	"context"
	"encoding/json"
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

type Storage interface {
	GetURLShort(originalURL string, userID string) string
	GetUserURLS(userID string) []URLPair
	Ping(ctx context.Context) error
	GetOriginalURL(shortURL string) (string, error)
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