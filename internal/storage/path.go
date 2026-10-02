package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

// DefaultDataDir returns the platform-standard user application data directory.
func DefaultDataDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		home, hErr := os.UserHomeDir()
		if hErr != nil {
			return "", fmt.Errorf("could not determine user data dir: %w", err)
		}
		configDir = filepath.Join(home, ".config")
	}

	appDir := filepath.Join(configDir, "acctg-practice")
	if err := os.MkdirAll(appDir, 0700); err != nil {
		return "", fmt.Errorf("could not create application data directory %s: %w", appDir, err)
	}
	return appDir, nil
}

// DefaultDBPath returns the full path to the default SQLite database.
func DefaultDBPath() (string, error) {
	dir, err := DefaultDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "acctg_practice.db"), nil
}
