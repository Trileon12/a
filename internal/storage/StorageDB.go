package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"io/ioutil"
	"os"
)

type StorageDB struct {
	Conf *Config
	URLs URLsType
	DB   *sql.DB
}

func (s *StorageDB) Ping(ctx context.Context) error {
	return s.DB.PingContext(ctx)
}

func (s *StorageDB) Close() {
	if s.DB != nil {
		s.DB.Close()
	}
}

func NewStorageDB(conf *Config) *StorageDB {
	s := &StorageDB{
		Conf: conf,
		DB:   nil,
	}

	var err error
	if conf.IsDBDefined() {
		s.DB, err = sql.Open("pgx", conf.DBAddress)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		s.Migrate()
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

func (s *StorageDB) GetURLsShort(originalURL []ShortURLItemRequest, userID string, host string) ([]ShortURLItemResponse, error) {

	res := make([]ShortURLItemResponse, len(originalURL))
	var errd error
	for i := range originalURL {
		shortURL, err := s.GetURLShort(originalURL[i].OriginalURL, userID)
		if errors.Is(err, ErrDuplicateOriginalURL) {
			errd = ErrDuplicateOriginalURL
		}
		res[i] = ShortURLItemResponse{
			ShortURL:      host + shortURL,
			CorrelationID: originalURL[i].CorrelationID,
		}
	}
	return res, errd
}

// GetURLShort Create and return short url for given original URL. Return the same short url for the same orginal URL
func (s *StorageDB) GetURLShort(originalURL string, userID string) (string, error) {

	shortURL := RandString(s.Conf.MaxLength)
	rows := s.DB.QueryRow("insert into urls (\"UserId\", \"OriginalURL\", \"ShortURL\")"+
		" values ($1, $2, $3) ON CONFLICT (\"OriginalURL\")  DO UPDATE SET \"id\"=EXCLUDED.\"id\"  RETURNING \"ShortURL\"", userID, originalURL, shortURL)

	var res sql.NullString
	err := rows.Scan(&res)
	if err != nil {
		return "", err
	}

	if res.String != shortURL {
		return res.String, ErrDuplicateOriginalURL
	} else {
		return shortURL, nil
	}

}

func (s *StorageDB) GetUserURLS(userID string) []URLPair {

	rows, err := s.DB.Query("Select \"OriginalURL\",\"ShortURL\" from urls where \"UserId\"=$1", userID)
	if err != nil {
		return nil
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()

	res := make([]URLPair, 0)
	for rows.Next() { // ?????????????????? ???? ???????? ??????????????
		var url URLPair
		err = rows.Scan(&url.OriginalURL, &url.ShortURL)
		if err != nil {

			return nil
		}
		res = append(res, url)
	}
	return res

}

// GetOriginalURL Func returns original url by short url
func (s *StorageDB) GetOriginalURL(shortURL string) (string, error) {

	rows := s.DB.QueryRow("Select \"OriginalURL\" from urls where \"ShortURL\"=$1", shortURL)

	var res sql.NullString
	err := rows.Scan(&res)
	if err != nil {
		return "", err
	}
	if res.Valid {
		return res.String, nil
	} else {
		return "", errors.New("URL ???? ????????????")
	}
}

func (s *StorageDB) Migrate() {

	var mgratescript = "create table if not exists  urls (" +
		"ID            serial" +
		"        constraint id_pk" +
		"            primary key," +
		"    \"UserId\"      text," +
		"    \"OriginalURL\" text," +
		"    \"ShortURL\"    text" +
		");" +
		"create unique index if not exists  id__index" +
		"    on urls (ID);" +
		"create unique index  if not exists  originalurl_index" +
		"    on urls (\"OriginalURL\");" +
		"create index if not exists  shorturl_index" +
		"    on urls (\"ShortURL\");" +
		"create index if not exists  \"userId__index\"" +
		"    on urls (\"UserId\");"
	s.DB.Exec(mgratescript)
}

func (s *StorageDB) SaveData() {

}
