package storage

import (
	"os"
	"path"
	"strings"

	"github.com/loghinalexandru/anchor/internal/config"
)

type Kind int

const (
	Local Kind = iota
	Git
)

type Storer interface {
	Init(args ...string) error
	Store(msg string) error
}

func New(k Kind, path string) Storer {
	switch k {
	case Git:
		storer, err := newGitStorage(path)
		if err != nil {
			panic(err)
		}

		return storer
	default:
		return newLocalStorage(path)
	}
}

func Parse(s string) Kind {
	switch strings.ToLower(s) {
	case "git":
		return Git
	default:
		return Local
	}
}

func Path() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("Cannot open storage directory")
	}

	return path.Join(home, config.StdDirName)
}
