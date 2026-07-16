package storage

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func New(dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	s := &Storage{db: db}
	if err := s.runMigrations(); err != nil {
		return nil, fmt.Errorf("migrations: %w", err)
	}

	slog.Info("storage initialized", "path", dbPath)
	return s, nil
}

func (s *Storage) runMigrations() error {
	if _, err := s.db.Exec(migrationInit); err != nil {
		return fmt.Errorf("exec init migration: %w", err)
	}
	slog.Info("migration applied", "file", "init")
	return nil
}

func (s *Storage) DB() *sql.DB {
	return s.db
}

func (s *Storage) Close() error {
	return s.db.Close()
}
