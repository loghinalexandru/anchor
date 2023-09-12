package storage

import (
	"errors"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func CloneWithSSH(path string, remote string) error {
	auth, err := ssh.NewSSHAgentAuth("git")
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

	tree, _ := repo.Worktree()
	status, err := tree.Status()
	if err != nil || status.IsClean() {
		return err
	}

	auth, err := ssh.NewSSHAgentAuth("git")
	if err != nil {
		return err
	}

	err = tree.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       auth,
	})

	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return err
	}

	_, err = tree.Add(".")
	if err != nil {
		return err
	}

	_, err = tree.Commit("Sync bookmarks", &git.CommitOptions{})
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
