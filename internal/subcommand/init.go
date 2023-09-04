package subcommand

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/peterbourgon/ff/v4"
)

type initCmd struct {
	command *ff.Command
	flags   *ff.CoreFlags
}

var (
	ErrHomeDir   = errors.New("could not open home directory")
	ErrCreateDir = errors.New("could not create anchor directory")
)

func RegisterInit(root *ff.Command, rootFlags *ff.CoreFlags) error {
	var cmd initCmd
	cmd.flags = ff.NewFlags("init").SetParent(rootFlags)
	_ = cmd.flags.Bool('r', "repository", false, "used in order to init a git repository")

	cmd.command = &ff.Command{
		Name:      "init",
		Usage:     "init",
		ShortHelp: "initialize a new empty home for anchor",
		Flags:     cmd.flags,
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

	root.Subcommands = append(root.Subcommands, cmd.command)

	return nil
}

func (c *initCmd) handle(res chan<- error) {
	defer close(res)

	home, err := os.UserHomeDir()

	if err != nil {
		res <- fmt.Errorf("%w with base: '%w'", ErrHomeDir, err)
	}

	if _, err := os.Stat(home + "/.anchor"); os.IsNotExist(err) {
		err = os.Mkdir(home+"/.anchor", os.ModeDir)

		if err != nil {
			res <- fmt.Errorf("%w with base: '%w'", ErrCreateDir, err)
		}
	}
}
