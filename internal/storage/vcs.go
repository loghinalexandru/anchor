package storage

import (
	"errors"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

const (
	stdUser   = "git"
	msgCommit = "Sync bookmarks"
)

func CloneWithSSH(path string, remote string) error {
	auth, err := ssh.NewSSHAgentAuth(stdUser)
	if err != nil {
		return err
	}

	_, err = git.PlainClone(path, false, &git.CloneOptions{
		URL:  remote,
		Auth: auth,
	})
	if err != nil {
		return err
	}

	return nil
}

func PushWithSSH(path string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	auth, err := ssh.NewSSHAgentAuth(stdUser)
	if err != nil {
		return err
	}

	tree, _ := repo.Worktree()
	err = tree.Pull(&git.PullOptions{
		RemoteName: git.DefaultRemoteName,
		Auth:       auth,
	})

	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return err
	}

	status, err := tree.Status()
	if err != nil || status.IsClean() {
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
		Auth:       auth,
	})

	if err != nil {
		return err
	}

	return nil
}
