package storage

import (
	"errors"
)

var URLs map[string]string
var URLsRevers map[string]string

func init() {
	URLs = make(map[string]string)

}

const shortUrlLen = 6

//Create and return short url for given original URL. Return the same short url for the same orginal Url
func GetUrlShort(originalURL string) string {

	if shortUrl, isExists := URLsRevers[originalURL]; isExists {
		return shortUrl
	}

	shortUrl := getUnicUrl()
	URLs[shortUrl] = originalURL

	return shortUrl
}

//Func returns original url by short url
func GetOriginalUrl(shortURL string) (string, error) {

	if originalUrl, isExists := URLs[shortURL]; isExists {
		return originalUrl, nil
	} else {
		return "", errors.New("URL не найден")
	}
}

func getUnicUrl() string {

	found := false
	shortUrl := RandString(shortUrlLen)

	for _, found = URLs[shortUrl]; found; {
		shortUrl = RandString(shortUrlLen)
	}
	return shortUrl
}
