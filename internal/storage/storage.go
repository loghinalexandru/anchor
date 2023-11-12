package storage

type Storer interface {
	Init(remote string) error
	Store() error
}

func New(kind string) (Storer, error) {
	switch kind {
	case "git":
		return NewGitStorage()
	default:
		return NewLocalStorage(), nil
	}
}
