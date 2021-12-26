package storage

import (
	"errors"
)

var URLs map[string]string
var URLsRevers map[string]string

func init() {
	URLs = make(map[string]string)

}

const shortURLLen = 6

//Create and return short url for given original URL. Return the same short url for the same orginal URL
func GetURLShort(originalURL string) string {

	if shortURL, isExists := URLsRevers[originalURL]; isExists {
		return shortURL
	}

	shortURL := getUnicURL()
	URLs[shortURL] = originalURL

	return shortURL
}

//Func returns original url by short url
func GetOriginalURL(shortURL string) (string, error) {

	if originalURL, isExists := URLs[shortURL]; isExists {
		return originalURL, nil
	} else {
		return "", errors.New("URL не найден")
	}
}

func getUnicURL() string {

	found := false
	shortURL := RandString(shortURLLen)

	for _, found = URLs[shortURL]; found; {
		shortURL = RandString(shortURLLen)
	}
	return shortURL
}
