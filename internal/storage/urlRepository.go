package storage

import (
	"errors"
)

type Config struct {
	MaxLength     int
	HostShortURLs string
}

type Storage struct {
	conf       *Config
	URLs       map[string]string
	URLsRevers map[string]string
}

func New(conf *Config) *Storage {
	return &Storage{
		conf:       conf,
		URLs:       make(map[string]string),
		URLsRevers: make(map[string]string),
	}
}

//Create and return short url for given original URL. Return the same short url for the same orginal URL
func (s *Storage) GetURLShort(originalURL string) string {

	if shortURL, isExists := s.URLsRevers[originalURL]; isExists {
		return shortURL
	}

	shortURL := s.getUnicURL()
	s.URLs[shortURL] = originalURL

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
