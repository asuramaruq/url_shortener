package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/asuramaruq/url_shortener/internal/storage"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.NewStorage" // название функции для логов

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS url(
				id INTEGER PRIMARY KEY,
				alias TEXT NOT NULL UNIQUE,
				url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`

	stmt, err := db.Prepare(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// SaveURL

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	insertQuery := `
		INSERT INTO url(url, alias) values(?, ?)
	`

	stmt, err := s.db.Prepare(insertQuery)
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

// GetURL

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	getQuery := `
		SELECT url FROM url WHERE alias = ?
	`

	stmt, err := s.db.Prepare(getQuery)
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL)

	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}

	if err != nil {
		return "", fmt.Errorf("%s: query execution: %w", op, err)
	}

	return resURL, nil
}

// DeleteURL

func (s *Storage) DeleteURL(alias string) (bool, error) {
	const op = "storage.sqlite.DeleteURL"

	deleteQuery := `
		DELETE FROM url WHERE alias = ?
	`

	stmt, err := s.db.Prepare(deleteQuery)
	if err != nil {
		return false, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	res, err := stmt.Exec(alias)
	if err != nil {
		return false, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("%s: rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return false, storage.ErrURLNotFound
	}

	return true, nil
}
