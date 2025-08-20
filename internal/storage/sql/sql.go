package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/internal/storage"

	"github.com/go-sql-driver/mysql"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(dsn string) (*Storage, error) {
	const op = "storage.mysql.New"
	dsn += "url_shortener"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt := `
    CREATE TABLE IF NOT EXISTS url (
        id BIGINT AUTO_INCREMENT PRIMARY KEY,
        alias VARCHAR(255) NOT NULL UNIQUE,
        url TEXT NOT NULL,
        INDEX idx_alias (alias)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
    `

	_, err = db.Exec(stmt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURl(alias string, urlToSave string) (int64, error) {
	const op = "storage.mysql.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url (alias, url) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(alias, urlToSave)
	if err != nil {

		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil

}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.mysql.GetURL"
	var url string
	stmt, err := s.db.Prepare("SELECT url From url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = stmt.QueryRow(alias).Scan(&url)

	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
	}
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}

/*func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.mysql.DeleteURL"
}*/
