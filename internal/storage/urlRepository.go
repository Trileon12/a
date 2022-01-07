package storage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	MaxLength int
	FilePath  string `env:"FILE_STORAGE_PATH"`
}

type URLsType map[string]string

type Storage struct {
	conf *Config
	URLs URLsType
}

func New(conf *Config) *Storage {
	s := &Storage{
		conf: conf,
		URLs: make(URLsType),
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

func ExtractJSONURLData(fileName string, s *URLsType) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	_ = json.Unmarshal(byteValue, &s)
}

func (s *Storage) SaveData() {

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
func (s *Storage) GetURLShort(originalURL string) string {

	shortURL := s.getUnicURL()
	s.URLs[shortURL] = originalURL

	s.SaveData()
	return shortURL
}

//Func returns original url by short url
func (s *Storage) GetOriginalURL(shortURL string) (string, error) {

	if originalURL, isExists := s.URLs[shortURL]; isExists {
		return originalURL, nil
	} else {
		return "", errors.New("URL не найден")
	}
}

func (s *Storage) getUnicURL() string {

	found := false
	shortURL := RandString(s.conf.MaxLength)

	for _, found = s.URLs[shortURL]; found; {
		shortURL = RandString(s.conf.MaxLength)
	}
	return shortURL
}
