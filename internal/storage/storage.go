package storage

import "strings"

type Kind int

const (
	Local Kind = iota
	Git
)

type Storer interface {
	Init(args ...string) error
	Store() error
}

func New(k Kind) Storer {
	switch k {
	case Git:
		storer, err := newGitStorage()
		if err != nil {
			panic(err)
		}
		return storer
	default:
		return newLocalStorage()
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
