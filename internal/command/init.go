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
	Command  ff.Command
	repoFlag bool
}

func NewInit(rootFlags *ff.FlagSet) *initCmd {
	cmd := initCmd{}

	flags := ff.NewFlagSet("init").SetParent(rootFlags)
	_ = flags.BoolVar(&cmd.repoFlag, 'r', "repository", "used in order to init a git repository")

	cmd.Command = ff.Command{
		Name:      "init",
		Usage:     "init",
		ShortHelp: "init a new empty home for anchor",
		Flags:     flags,
		Exec:      handlerMiddleware(cmd.handle),
	}

	return &cmd
}

func (init *initCmd) handle(_ context.Context, args []string) error {
	dir, err := config.RootDir()
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
		err = os.Mkdir(dir, config.StdFileMode)
		if err != nil {
			return err
		}
	}

	return nil
}
