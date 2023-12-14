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

const (
	msgUpdateFailed = "Failed pulling latest changes. Continue operation?"
)

type storerContextKey struct{}

type Updater interface {
	Update() error
}

type rootCmd struct {
	storage string
	cmd     *ff.Command
}

func newRoot() *rootCmd {
	root := &rootCmd{}

	rootFlags := ff.NewFlagSet("anchor")
	rootFlags.StringVar(&root.storage, 's', "storage", "local", "Set storage type")

	root.cmd = &ff.Command{
		Name:  "anchor",
		Usage: "anchor [FLAGS] <SUBCOMMAND>",
		Flags: rootFlags,
	}

	root.cmd.Subcommands = []*ff.Command{
		(&getCmd{}).manifest(rootFlags),
		(&createCmd{}).manifest(rootFlags),
		(&treeCmd{}).manifest(rootFlags),
		(&importCmd{}).manifest(rootFlags),
		(&deleteCmd{}).manifest(rootFlags),
		(&initCmd{}).manifest(rootFlags),
		(&syncCmd{}).manifest(rootFlags),
	}

	return root
}

func (root *rootCmd) handle(ctx context.Context, args []string) error {
	err := root.cmd.Parse(args,
		ff.WithConfigFile(config.FilePath()),
		ff.WithConfigFileParser(ffyaml.Parse),
		ff.WithConfigAllowMissingFile())

	if err != nil {
		return err
	}

	storer := storage.New(storage.Parse(root.storage))
	for _, c := range root.cmd.Subcommands {
		if updater, ok := storer.(Updater); ok {
			c.Exec = updaterMiddleware(c.Exec, updater)
		}

		c.Exec = contextHandlerMiddleware(c.Exec)
	}

	return root.cmd.Run(context.WithValue(ctx, storerContextKey{}, storer))
}

type handlerFunc func(ctx context.Context, args []string) error

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
