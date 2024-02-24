package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

func (s *Storage) GetFile(ctx context.Context, filename string) ([]byte, error) {
	const op = "storage.postgres.GetFile"

	query := "SELECT file_data FROM files WHERE filename = $1"
	row := s.db.QueryRowContext(ctx, query, filename)

	var fileData []byte
	err := row.Scan(&fileData)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return fileData, nil
}

func (s *Storage) Upload(filename string, fileData []byte) error {
	const op = "storage.postgres.Upload"

	query := "INSERT INTO files(filename, file_data) VALUES($1, $2)"
	_, err := s.db.Exec(query, filename, fileData)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
