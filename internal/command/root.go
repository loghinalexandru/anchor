package command

import (
	"context"

	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/loghinalexandru/anchor/internal/output"
	"github.com/loghinalexandru/anchor/internal/storage"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffyaml"
)

const (
	rootName        = "anchor"
	msgUpdateFailed = "Failed pulling latest changes. Continue operation?"
)

type Updater interface {
	Update() error
}

type rootContext struct {
	context.Context
	storer storage.Storer
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
		Name:  rootName,
		Usage: "anchor [FLAGS] <SUBCOMMAND>",
		Flags: rootFlags,
	}

	root.cmd.Subcommands = []*ff.Command{
		(&initCmd{}).manifest(rootFlags),
		(&viewCmd{}).manifest(rootFlags),
		(&addCmd{}).manifest(rootFlags),
		(&deleteCmd{}).manifest(rootFlags),
		(&treeCmd{}).manifest(rootFlags),
		(&syncCmd{}).manifest(rootFlags),
		(&importCmd{}).manifest(rootFlags),
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
		switch c.Name {
		// Skip updateMiddleware for init command.
		case initName:
			c.Exec = contextMiddleware(c.Exec)
		default:
			c.Exec = contextMiddleware(updaterMiddleware(c.Exec, storer))
		}
	}

	cmdCtx := rootContext{
		Context: ctx,
		storer:  storer,
	}

	return root.cmd.Run(cmdCtx)
}

type handlerFunc func(ctx context.Context, args []string) error

func updaterMiddleware(next handlerFunc, storer storage.Storer) handlerFunc {
	updater, ok := storer.(Updater)
	if !ok {
		return next
	}

	return func(ctx context.Context, args []string) error {
		err := updater.Update()
		if err != nil {
			if ok := output.Confirm(msgUpdateFailed); !ok {
				return err
			}
		}

		return next(ctx, args)
	}
}

func contextMiddleware(next handlerFunc) handlerFunc {
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
