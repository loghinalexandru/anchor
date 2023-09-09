package vcs

import (
	"errors"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func PlainInitWithRemote(path string, remote string) error {
	repo, err := git.PlainInit(path, false)
	if err != nil {
		return err
	}

	// Validate remote url
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name:  "origin",
		URLs:  []string{remote},
		Fetch: []config.RefSpec{config.RefSpec("+refs/heads/*:refs/remotes/origin/*")},
	})

	if err != nil {
		return err
	}

	return nil
}

func PullWithSSH(path string) error {
	repo, err := git.PlainOpen(path)

	if err != nil {
		return err
	}

	tree, err := repo.Worktree()
	if err != nil {
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
	if err != nil {
		return err
	}

	if status.IsClean() {
		return nil
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
		RemoteName: "origin",
		Auth:       auth,
	})

	if err != nil {
		return err
	}

	return nil
}
