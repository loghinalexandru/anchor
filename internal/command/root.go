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
	storage string
}

type handlerFunc func(ctx context.Context, args []string) error

func newRoot() *rootCmd {
	var root rootCmd

	rootFlags := ff.NewFlagSet("anchor")
	_ = rootFlags.StringVar(&root.storage, 's', "storage", "local", "Set storage type")

	root.cmd = ff.Command{
		Name:  "anchor",
		Usage: "anchor [FLAGS] <SUBCOMMAND>",
		Flags: rootFlags,
	}

	return &root
}

func (root *rootCmd) bootstrap(args []string) error {
	flags := root.cmd.Flags.(*ff.FlagSet)
	initCmd := newInit(flags)
	syncCmd := newSync(flags)

	root.cmd.Subcommands = append(root.cmd.Subcommands, &newCreate(flags).command)
	root.cmd.Subcommands = append(root.cmd.Subcommands, &initCmd.command)
	root.cmd.Subcommands = append(root.cmd.Subcommands, &newGet(flags).command)
	root.cmd.Subcommands = append(root.cmd.Subcommands, &newDelete(flags).command)
	root.cmd.Subcommands = append(root.cmd.Subcommands, &syncCmd.command)
	root.cmd.Subcommands = append(root.cmd.Subcommands, (*ff.Command)(newImport(flags)))
	root.cmd.Subcommands = append(root.cmd.Subcommands, (*ff.Command)(newTree(flags)))

	err := root.cmd.Parse(args, ff.WithConfigFile(config.FilePath()), ff.WithConfigFileParser(ffyaml.Parse))
	if err != nil {
		return err
	}

	store, err := storage.New(root.storage)
	if err != nil {
		return err
	}

	initCmd.withStorage(store)
	syncCmd.withStorage(store)

	for _, c := range root.cmd.Subcommands {
		if updater, ok := store.(Updater); ok {
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
