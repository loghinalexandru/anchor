package command

import (
	"context"

	"github.com/peterbourgon/ff/v4"
)

func newExec() *ff.Command {
	rootFlags := ff.NewFlagSet("anchor")
	rootCmd := &ff.Command{
		Name:  "anchor",
		Usage: "anchor [FLAGS] <SUBCOMMAND>",
		Flags: rootFlags,
	}

	rootCmd.Subcommands = append(rootCmd.Subcommands, &newCreate(rootFlags).command)
	rootCmd.Subcommands = append(rootCmd.Subcommands, &newInit(rootFlags).command)
	rootCmd.Subcommands = append(rootCmd.Subcommands, &newGet(rootFlags).command)
	rootCmd.Subcommands = append(rootCmd.Subcommands, &newDelete(rootFlags).command)
	rootCmd.Subcommands = append(rootCmd.Subcommands, (*ff.Command)(newSync(rootFlags)))
	rootCmd.Subcommands = append(rootCmd.Subcommands, (*ff.Command)(newImport(rootFlags)))
	rootCmd.Subcommands = append(rootCmd.Subcommands, (*ff.Command)(newTree(rootFlags)))

	return rootCmd
}

type handlerFunc func(ctx context.Context, args []string) error

func handlerMiddleware(next handlerFunc) handlerFunc {
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