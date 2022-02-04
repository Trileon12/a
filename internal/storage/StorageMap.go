package storage

import (
	"context"
	"encoding/json"
	"errors"
	_ "github.com/jackc/pgx/v4/stdlib"
	"io/ioutil"
	"log"
	"os"
)

type UserURLs map[string][]URLPair

type StorageMap struct {
	conf     *Config
	URLs     URLsType
	UserURLs UserURLs
}

func NewStorageMap(conf *Config) Storage {
	s := &StorageMap{
		conf:     conf,
		URLs:     make(URLsType),
		UserURLs: make(UserURLs),
	}

	if len(conf.FilePath) > 0 {
		ExtractJSONURLData(conf.FilePath, &s.URLs)
		flag1 := os.O_WRONLY | os.O_CREATE | os.O_APPEND
		jsonFile, err := os.OpenFile(conf.FilePath, flag1, 0777)
		if err != nil {
			return s
		}
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)
		_ = json.Unmarshal(byteValue, &s)

	}

	return s
}

func (s *StorageMap) GetURLsShort(originalURL []ShortURLItemRequest, userID string, host string) []ShortURLItemResponse {

	res := make([]ShortURLItemResponse, len(originalURL))
	for i, _ := range originalURL {
		shortURL := s.GetURLShort(originalURL[i].OriginalURL, userID)
		res = append(res, ShortURLItemResponse{
			ShortURL:      host + shortURL,
			CorrelationID: originalURL[i].CorrelationID,
		})
	}
	return res
}

func (s *StorageMap) Ping(ctx context.Context) error {
	return nil
}
func (s *StorageMap) Close() {

}
func (s *StorageMap) SaveData() {

	if len(s.conf.FilePath) > 0 {
		jsonString, err := json.Marshal(s.URLs)
		if err != nil {
			log.Fatal(err)
			return
		}
		_ = ioutil.WriteFile(s.conf.FilePath, jsonString, 0644)

	}
}

//Create and return short url for given original URL. Return the same short url for the same orginal URL
func (s *StorageMap) GetURLShort(originalURL string, userID string) string {

	shortURL := s.getUnicURL()
	s.URLs[shortURL] = originalURL

	newPair := URLPair{
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}

	s.UserURLs[userID] = append(s.UserURLs[userID], newPair)
	s.SaveData()
	return shortURL
}

func (s *StorageMap) GetUserURLS(userID string) []URLPair {
	return s.UserURLs[userID]

}

// GetOriginalURL Func returns original url by short url
func (s *StorageMap) GetOriginalURL(shortURL string) (string, error) {

	if originalURL, isExists := s.URLs[shortURL]; isExists {
		return originalURL, nil
	} else {
		return "", errors.New("URL не найден")
	}
}

func (s *StorageMap) getUnicURL() string {

	found := false
	shortURL := RandString(s.conf.MaxLength)

	for _, found = s.URLs[shortURL]; found; {
		shortURL = RandString(s.conf.MaxLength)
	}
	return shortURL
}
