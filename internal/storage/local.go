package storage

import (
	"os"

	"github.com/loghinalexandru/anchor/internal/config"
)

type localStorage struct {
	path string
}

func newLocalStorage() *localStorage {
	return &localStorage{
		path: config.RootDir(),
	}
}

func (*localStorage) Init(_ string) error {
	dir := config.RootDir()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func (*localStorage) Store() error {
	return nil
}
