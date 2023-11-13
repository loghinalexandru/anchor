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

type handlerFunc func(ctx context.Context, args []string) error

func newRoot(args []string) (*ff.Command, error) {
	var storageKind string

	rootFlags := ff.NewFlagSet("anchor")
	rootFlags.StringVar(&storageKind, 's', "storage", "local", "Set storage type")

	cmd := &ff.Command{
		Name:  "anchor",
		Usage: "anchor [FLAGS] <SUBCOMMAND>",
		Flags: rootFlags,
	}

	initialize := newInit()
	sync := newSync()

	cmd.Subcommands = []*ff.Command{
		(&getCmd{}).manifest(rootFlags),
		(&createCmd{}).manifest(rootFlags),
		(&treeCmd{}).manifest(rootFlags),
		(&importCmd{}).manifest(rootFlags),
		(&deleteCmd{}).manifest(rootFlags),
		initialize.manifest(rootFlags),
		sync.manifest(rootFlags),
	}

	err := cmd.Parse(args,
		ff.WithConfigFile(config.FilePath()),
		ff.WithConfigFileParser(ffyaml.Parse),
		ff.WithConfigAllowMissingFile())
	if err != nil {
		return nil, err
	}

	storer := storage.New(storage.Parse(storageKind))

	// Assign storer implementation determined after parsing
	initialize.withStorage(storer)
	sync.withStorage(storer)

	for _, c := range cmd.Subcommands {
		if updater, ok := storer.(Updater); ok {
			c.Exec = updaterMiddleware(c.Exec, updater)
		}

		c.Exec = contextHandlerMiddleware(c.Exec)
	}

	return cmd, nil
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
