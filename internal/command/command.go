package command

import (
	"context"

	"github.com/peterbourgon/ff/v4"
)

func NewExec() *ff.Command {
	rootFlags := ff.NewFlagSet("anchor")
	rootCmd := &ff.Command{
		Name:  "anchor",
		Usage: "anchor [FLAGS] <SUBCOMMAND>",
		Flags: rootFlags,
	}

	rootCmd.Subcommands = append(rootCmd.Subcommands, &NewCreate(rootFlags).Command)
	rootCmd.Subcommands = append(rootCmd.Subcommands, &NewInit(rootFlags).Command)
	rootCmd.Subcommands = append(rootCmd.Subcommands, &NewGet(rootFlags).Command)
	rootCmd.Subcommands = append(rootCmd.Subcommands, &NewDelete(rootFlags).Command)
	rootCmd.Subcommands = append(rootCmd.Subcommands, (*ff.Command)(NewSync(rootFlags)))
	rootCmd.Subcommands = append(rootCmd.Subcommands, (*ff.Command)(NewImport(rootFlags)))
	rootCmd.Subcommands = append(rootCmd.Subcommands, (*ff.Command)(NewTree(rootFlags)))

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
