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

type command interface {
	def() *ff.Command
	handle(context.Context, []string) error
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
	}, "Set storage type (default: local)")

	root.cmd = ff.Command{
		Name:  "anchor",
		Usage: "anchor [FLAGS] <SUBCOMMAND>",
		Flags: rootFlags,
	}

	return &root
}

func (root *rootCmd) bootstrap(args []string) error {
	flags := root.cmd.Flags.(*ff.FlagSet)

	subcommands := []command{
		newCreate(flags),
		newInit(flags),
		newGet(flags),
		newDelete(flags),
		newSync(flags),
		newImport(flags),
		newTree(flags),
	}

	root.cmd.Subcommands = make([]*ff.Command, len(subcommands))
	for i, cm := range subcommands {
		root.cmd.Subcommands[i] = cm.def()
	}

	err := root.cmd.Parse(args,
		ff.WithConfigFile(config.FilePath()),
		ff.WithConfigFileParser(ffyaml.Parse),
		ff.WithConfigAllowMissingFile())

	if err != nil {
		return err
	}

	storer := storage.New(root.storage)
	for _, c := range subcommands {
		if setter, ok := c.(interface{ withStorage(storage.Storer) }); ok {
			setter.withStorage(storer)
		}

		if updater, ok := storer.(Updater); ok {
			c.def().Exec = updaterMiddleware(c.def().Exec, updater)
		}

		c.def().Exec = contextHandlerMiddleware(c.def().Exec)
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
