package subcommand

import (
	"context"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/peterbourgon/ff/v4"
)

type syncCmd ff.Command

func RegisterSync(root *ff.Command, rootFlags *ff.CoreFlags) {
	var cmd *syncCmd
	flags := ff.NewFlags("sync").SetParent(rootFlags)

	cmd = &syncCmd{
		Name:      "sync",
		Usage:     "sync",
		ShortHelp: "sync changes with configured remote",
		Flags:     flags,
		Exec: func(ctx context.Context, args []string) error {
			res := make(chan error, 1)
			go cmd.handle(res)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-res:
				return err
			}
		},
	}

	root.Subcommands = append(root.Subcommands, (*ff.Command)(cmd))
}

func (c *syncCmd) handle(res chan<- error) {
	dir, _ := c.Flags.GetFlag("root-dir")
	home, err := os.UserHomeDir()

	if err != nil {
		res <- err
		return
	}

	path := filepath.Join(home, dir.GetValue())
	repo, err := git.PlainOpen(path)

	if err != nil {
		res <- err
		return
	}

	tree, _ := repo.Worktree()
	status, err := tree.Status()
	if err != nil {
		res <- err
		return
	}

	if status.IsClean() {
		close(res)
		return
	}

	auth, err := ssh.NewSSHAgentAuth("git")
	if err != nil {
		res <- err
		return
	}

	_ = tree.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       auth,
	})

	_, err = tree.Add(".")
	if err != nil {
		res <- err
		return
	}

	_, err = tree.Commit("Sync bookmarks", &git.CommitOptions{})
	if err != nil {
		res <- err
		return
	}

	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
	})

	if err != nil {
		res <- err
	}

	close(res)
}
