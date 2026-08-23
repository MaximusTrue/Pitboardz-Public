package logging

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	pitboardzzDirectoryName = "PitBoardzz"
	dataDirectoryName       = "data"
	databaseFileName        = "pitboardz.db"
)

// Store owns PitBoardzz's local logging database connection.
type Store struct {
	database     *sql.DB
	databasePath string
}

// Open creates the logging data directory and opens the SQLite database.
func Open(savePath string) (*Store, error) {
	if savePath == "" {
		return nil, errors.New("save path is empty")
	}

	dataPath := filepath.Join(savePath, pitboardzzDirectoryName, dataDirectoryName)
	if err := os.MkdirAll(dataPath, 0o755); err != nil {
		return nil, fmt.Errorf("create logging data directory: %w", err)
	}

	databasePath := filepath.Join(dataPath, databaseFileName)
	database, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, fmt.Errorf("open logging database: %w", err)
	}
	database.SetMaxOpenConns(1)

	if err := database.Ping(); err != nil {
		_ = database.Close()
		return nil, fmt.Errorf("connect to logging database: %w", err)
	}

	return &Store{database: database, databasePath: databasePath}, nil
}

// Path returns the full path to the SQLite database file.
func (store *Store) Path() string {
	if store == nil {
		return ""
	}
	return store.databasePath
}

// Close releases the SQLite database connection.
func (store *Store) Close() error {
	if store == nil || store.database == nil {
		return nil
	}
	return store.database.Close()
}
