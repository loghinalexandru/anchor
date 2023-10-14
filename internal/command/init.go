package command

import (
	"context"
	"errors"
	"os"

	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/loghinalexandru/anchor/internal/storage"
	"github.com/peterbourgon/ff/v4"
)

var (
	ErrInvalidURL = errors.New("not a valid URL")
)

type initCmd struct {
	command  ff.Command
	repoFlag bool
}

func newInit(rootFlags *ff.FlagSet) *initCmd {
	cmd := initCmd{}

	flags := ff.NewFlagSet("init").SetParent(rootFlags)
	_ = flags.BoolVar(&cmd.repoFlag, 'r', "repository", "used in order to init a git repository")

	cmd.command = ff.Command{
		Name:      "init",
		Usage:     "init",
		ShortHelp: "init a new empty home for anchor",
		Flags:     flags,
		Exec:      cmd.handle,
	}

	return &cmd
}

func (init *initCmd) handle(_ context.Context, args []string) error {
	if init.repoFlag {
		if len(args) == 0 {
			return ErrInvalidURL
		}

		s, _ := storage.NewGitStorage()
		err := s.CloneWithSSH(args[0])
		if err != nil {
			return err
		}

		return nil
	}

	dir := config.RootDir()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, config.StdFileMode)
		if err != nil {
			return err
		}
	}

	return nil
}
