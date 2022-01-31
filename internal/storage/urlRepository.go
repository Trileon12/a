package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	MaxLength int    `env:"MaxLength" envDefault:"6"`
	FilePath  string `env:"FILE_STORAGE_PATH"`
	DBAddress string `env:"DATABASE_DSN" envDefault:"postgres://postgres:myPassword@localhost:5432/dbvideo"`
}

type URLsType map[string]string
type URLPair struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type UserURLs map[string][]URLPair

type Storage struct {
	conf     *Config
	URLs     URLsType
	UserURLs UserURLs
}

type StoragePG struct {
	conf *Config
	DB   *pgx.Conn
}

func NewPG(conf *Config) *StoragePG {
	db, err := pgx.Connect(context.Background(), conf.DBAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return &StoragePG{
		conf: conf,
		DB:   db,
	}
}

func (s *StoragePG) Ping(ctx context.Context) error {
	return s.DB.Ping(ctx)

}

func (s *StoragePG) Close(ctx context.Context) {
	s.DB.Close(ctx)

}

func New(conf *Config) *Storage {
	s := &Storage{
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
func (s *Storage) GetURLShort(originalURL string, userID string) string {

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

func (s *Storage) GetUserURLS(userID string) []URLPair {
	return s.UserURLs[userID]

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
