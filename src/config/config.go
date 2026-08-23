// Package config owns the PitBoardzz INI file lifecycle.
package config

import (
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

// Store provides access to the runtime configuration file.
type Store struct {
	savePath string
}

func NewStore(savePath string) *Store {
	return &Store{savePath: savePath}
}

func (store *Store) Available() bool {
	return store != nil && store.savePath != ""
}

func (store *Store) EnsureDefault(contents string) (bool, error) {
	if !store.Available() {
		return false, nil
	}
	configPath := store.Path()
	if _, err := os.Stat(configPath); err == nil {
		return false, nil
	} else if !os.IsNotExist(err) {
		return false, err
	}
	return true, os.WriteFile(configPath, []byte(contents), 0o644)
}

func (store *Store) Load() (*ini.File, error) {
	if !store.Available() {
		return ini.Empty(), nil
	}
	return ini.Load(store.Path())
}

func (store *Store) Save(configuration *ini.File) error {
	if !store.Available() {
		return nil
	}
	return configuration.SaveTo(store.Path())
}

func (store *Store) Path() string {
	if !store.Available() {
		return ""
	}
	return filepath.Join(store.savePath, "pitboardz.ini")
}
