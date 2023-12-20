package storage

import (
	"errors"
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

func (l *localStorage) Init(_ ...string) error {
	if _, err := os.Stat(l.path); os.IsNotExist(err) {
		err = os.Mkdir(l.path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func (*localStorage) Diff() (string, error) {
	return "", errors.New("running on local storage type, command has no effect")
}

func (*localStorage) Store(_ string) error {
	return nil
}
