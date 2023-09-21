package types

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
)

type handlerFunc func(ctx context.Context, args []string) error

func Execute(args []string) error {
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

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	err := rootCmd.ParseAndRun(ctx, args)

	if errors.Is(err, ff.ErrHelp) || errors.Is(err, ff.ErrNoExec) {
		_, _ = fmt.Fprint(os.Stdout, ffhelp.Command(rootCmd))
		return nil
	}

	return err
}

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
