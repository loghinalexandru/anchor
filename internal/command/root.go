package command

import (
	"context"
	"os"

	"github.com/loghinalexandru/anchor/internal/output"
	"github.com/loghinalexandru/anchor/internal/storage"
	"github.com/peterbourgon/ff/v4"
)

var (
	msgUpdateFailed = "Failed pulling latest changes. Continue operation?"
)

type Updater interface {
	Update() error
}

func newRoot() *ff.Command {
	store, _ := storage.NewGitStorage()

	rootFlags := ff.NewFlagSet("anchor")
	root := &ff.Command{
		Name:  "anchor",
		Usage: "anchor [FLAGS] <SUBCOMMAND>",
		Flags: rootFlags,
	}

	root.Subcommands = append(root.Subcommands, &newCreate(rootFlags).command)
	root.Subcommands = append(root.Subcommands, &newInit(rootFlags).command)
	root.Subcommands = append(root.Subcommands, &newGet(rootFlags).command)
	root.Subcommands = append(root.Subcommands, &newDelete(rootFlags).command)
	root.Subcommands = append(root.Subcommands, &newSync(rootFlags, store).command)
	root.Subcommands = append(root.Subcommands, (*ff.Command)(newImport(rootFlags)))
	root.Subcommands = append(root.Subcommands, (*ff.Command)(newTree(rootFlags)))

	for _, c := range root.Subcommands {
		c.Exec = updaterMiddleware(contextHandlerMiddleware(c.Exec), store)
	}

	return root
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
