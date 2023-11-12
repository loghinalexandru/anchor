package storage

type Storer interface {
	Init(remote string) error
	Store() error
}

func New(kind string) (Storer, error) {
	switch kind {
	case "git":
		return newGitStorage()
	default:
		return newLocalStorage(), nil
	}
}
