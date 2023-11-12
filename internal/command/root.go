package command

import (
	"context"
	"os"

	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/loghinalexandru/anchor/internal/output"
	"github.com/loghinalexandru/anchor/internal/storage"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffyaml"
)

var (
	msgUpdateFailed = "Failed pulling latest changes. Continue operation?"
)

type Updater interface {
	Update() error
}

type rootCmd struct {
	cmd     ff.Command
	storage storage.Kind
}

type handlerFunc func(ctx context.Context, args []string) error

func newRoot() *rootCmd {
	var root rootCmd

	rootFlags := ff.NewFlagSet("anchor")
	rootFlags.Func('s', "storage", func(flag string) error {
		root.storage = storage.Parse(flag)
		return nil
	}, "Set storage type")

	root.cmd = ff.Command{
		Name:  "anchor",
		Usage: "anchor [FLAGS] <SUBCOMMAND>",
		Flags: rootFlags,
	}

	return &root
}

func (root *rootCmd) bootstrap(args []string) error {
	var storer storage.Storer

	flags := root.cmd.Flags.(*ff.FlagSet)
	root.cmd.Subcommands = append(root.cmd.Subcommands, &newCreate(flags).command)
	root.cmd.Subcommands = append(root.cmd.Subcommands, &newInit(flags, &storer).command)
	root.cmd.Subcommands = append(root.cmd.Subcommands, &newGet(flags).command)
	root.cmd.Subcommands = append(root.cmd.Subcommands, &newDelete(flags).command)
	root.cmd.Subcommands = append(root.cmd.Subcommands, &newSync(flags, &storer).command)
	root.cmd.Subcommands = append(root.cmd.Subcommands, (*ff.Command)(newImport(flags)))
	root.cmd.Subcommands = append(root.cmd.Subcommands, (*ff.Command)(newTree(flags)))

	err := root.cmd.Parse(args, ff.WithConfigFile(config.FilePath()), ff.WithConfigFileParser(ffyaml.Parse))
	if err != nil {
		return err
	}

	storer, err = storage.New(root.storage)
	if err != nil {
		return err
	}

	for _, c := range root.cmd.Subcommands {
		if updater, ok := storer.(Updater); ok {
			c.Exec = updaterMiddleware(c.Exec, updater)
		}

		c.Exec = contextHandlerMiddleware(c.Exec)
	}

	return nil
}

func updaterMiddleware(next handlerFunc, updater Updater) handlerFunc {
	return func(ctx context.Context, args []string) error {
		err := updater.Update()
		if err != nil {
			if ok := output.Confirmation(msgUpdateFailed, os.Stdin, os.Stdout); !ok {
				return err
			}
		}

		return next(ctx, args)
	}
}

func contextHandlerMiddleware(next handlerFunc) handlerFunc {
	return func(ctx context.Context, args []string) error {
		res := make(chan error, 1)

		go func(res chan<- error) {
			res <- next(ctx, args)
			close(res)
		}(res)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-res:
			return err
		}
	}
}
