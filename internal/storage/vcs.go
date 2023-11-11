package storage

import (
	"errors"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/loghinalexandru/anchor/internal/config"
)

const (
	stdUser   = "git"
	msgCommit = "Sync bookmarks"
)

type GitStorage struct {
	path string
	auth transport.AuthMethod
}

func NewGitStorage() (*GitStorage, error) {
	dir := config.RootDir()

	auth, err := ssh.NewSSHAgentAuth(stdUser)
	if err != nil {
		return nil, err
	}

	return &GitStorage{
		path: dir,
		auth: auth,
	}, nil
}

func (storage *GitStorage) CloneWithSSH(remote string) error {
	_, err := git.PlainClone(storage.path, false, &git.CloneOptions{
		URL:  remote,
		Auth: storage.auth,
	})

	return err
}

func (storage *GitStorage) Update() error {
	repo, err := git.PlainOpen(storage.path)
	if err != nil {
		return err
	}

	tree, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = tree.Pull(&git.PullOptions{
		RemoteName: git.DefaultRemoteName,
		Auth:       storage.auth,
	})

	if errors.Is(err, git.NoErrAlreadyUpToDate) {
		return nil
	}

	return err
}

func (storage *GitStorage) Store() error {
	repo, err := git.PlainOpen(storage.path)
	if err != nil {
		return err
	}

	tree, err := repo.Worktree()
	if err != nil {
		return err
	}

	_, err = tree.Add(".")
	if err != nil {
		return err
	}

	_, err = tree.Commit(msgCommit, &git.CommitOptions{})
	if err != nil {
		return err
	}

	err = repo.Push(&git.PushOptions{
		RemoteName: git.DefaultRemoteName,
		Auth:       storage.auth,
	})

	return err
}

func (storage *GitStorage) Diff() (string, error) {
	repo, err := git.PlainOpen(storage.path)
	if err != nil {
		return "", err
	}

	tree, err := repo.Worktree()
	if err != nil {
		return "", err
	}

	status, err := tree.Status()
	if err != nil {
		return "", err
	}

	return status.String(), nil
}
