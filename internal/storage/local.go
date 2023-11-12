package storage

import (
	"os"

	"github.com/loghinalexandru/anchor/internal/config"
)

type LocalStorage struct {
	path string
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		path: config.RootDir(),
	}
}

func (*LocalStorage) Init(_ string) error {
	dir := config.RootDir()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func (*LocalStorage) Update() error {
	return nil
}

func (*LocalStorage) Store() error {
	return nil
}
