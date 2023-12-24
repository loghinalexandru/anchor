package command

import (
	"context"
	"errors"
	"io/fs"
	"os"

	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/loghinalexandru/anchor/internal/output"
	"github.com/loghinalexandru/anchor/internal/storage"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffyaml"
)

const (
	rootName        = "anchor"
	rootUsage       = "anchor <SUBCOMMAND>"
	msgUpdateFailed = "Failed pulling latest changes. Continue operation?"
)

type Updater interface {
	Update() error
}

// rootContext is a context.Context wrapper for
// type safety and to avoid key-value pairs.
type rootContext struct {
	context.Context
	storer storage.Storer
}

type rootCmd struct {
	cmd *ff.Command
}

func newRoot() *rootCmd {
	root := &rootCmd{}

	rootFlags := ff.NewFlagSet("anchor")

	root.cmd = &ff.Command{
		Name:  rootName,
		Usage: rootUsage,
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
		(&versionCmd{}).manifest(rootFlags),
	}

	return root
}

func (root *rootCmd) handle(ctx context.Context, args []string) error {
	err := root.cmd.Parse(args)
	if err != nil {
		return err
	}

	fh, err := os.Open(config.FilePath())
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	defer fh.Close()

	// Config file might not exist, ignore errors if so.
	storer := storage.New(storage.Local)
	ffyaml.Parse(fh, func(key, value string) error {
		if key == config.StdStorageKey {
			storer = storage.New(storage.Parse(value))
		}

		return nil
	})

	// Add appropriate middleware for each subcommand
	for _, c := range root.cmd.Subcommands {
		switch c.Name {
		// Skip updateMiddleware for commands that
		// do not need to fetch from remote.
		case initName:
			c.Exec = contextMiddleware(c.Exec)
		case versionName:
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
