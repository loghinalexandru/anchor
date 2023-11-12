package storage

import "strings"

type Kind int

const (
	Local Kind = iota
	Git
)

type Storer interface {
	Init(remote string) error
	Store() error
}

func New(k Kind) (Storer, error) {
	switch k {
	case Git:
		return newGitStorage()
	default:
		return newLocalStorage(), nil
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
