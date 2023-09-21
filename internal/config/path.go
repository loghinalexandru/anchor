package config

import (
	"os"
	"path/filepath"
)

func RootDir() (string, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(home, StdDir), nil
}
