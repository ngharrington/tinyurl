package store

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteUrlStore struct {
	db *sql.DB
}

func (s *SqliteUrlStore) Store(url string) (int, error) {
	query := `
	INSERT INTO urls ( url ) VALUES ( ? )
	`
	res, err := s.db.Exec(query, url)
	if err != nil {
		return 0, err
	}
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s *SqliteUrlStore) GetById(id int) (string, error) {
	query := "SELECT url FROM urls WHERE id=?"
	row := s.db.QueryRow(query, id)
	var url string
	err := row.Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}

func NewSqliteUrlStore(path string) (*SqliteUrlStore, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return &SqliteUrlStore{}, err
	}
	return &SqliteUrlStore{db: db}, nil
}
