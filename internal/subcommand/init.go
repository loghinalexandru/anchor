package subcommand

import (
	"context"
	"errors"
	"os"

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

func RegisterInit(root *ff.Command, rootFlags *ff.FlagSet) {
	cmd := initCmd{}

	flags := ff.NewFlagSet("init").SetParent(rootFlags)
	_ = flags.BoolVar(&cmd.repoFlag, 'r', "repository", "used in order to init a git repository")

	cmd.command = ff.Command{
		Name:      "init",
		Usage:     "init",
		ShortHelp: "init a new empty home for anchor",
		Flags:     flags,
		Exec:      handlerMiddleware(cmd.handle),
	}

	root.Subcommands = append(root.Subcommands, &cmd.command)
}

func (init *initCmd) handle(_ context.Context, args []string) error {
	dir, err := rootDir()
	if err != nil {
		return err
	}

	if init.repoFlag {
		if len(args) == 0 {
			return err
		}

		err = storage.CloneWithSSH(dir, args[0])
		if err != nil {
			return err
		}

		return nil
	}

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, stdFileMode)
		if err != nil {
			return err
		}
	}

	return nil
}
