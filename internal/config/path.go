package config

import (
	"os"
	"path/filepath"
)

type pathFunc func() (string, error)

func RootDir(ff ...pathFunc) string {
	ff = append(ff, stdDir)
	for _, f := range ff {
		dir, err := f()
		if err == nil {
			return dir
		}
	}

	panic("Can not determine root dir")
}

func stdDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".anchor"), nil
}
