package subcommand

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/peterbourgon/ff/v4"
)

var (
	ErrInvalidURL = errors.New("not a valid url")
)

type createCmd ff.Command

func RegisterCreate(root *ff.Command, rootFlags *ff.CoreFlags) {
	var cmd *createCmd
	flags := ff.NewFlags("create").SetParent(rootFlags)
	_ = flags.String('l', "label", "", "add label in order of appearance")
	_ = flags.String('t', "title", "", "add custom title")

	cmd = &createCmd{
		Name:      "create",
		Usage:     "crate",
		ShortHelp: "add a bookmark with set labels",
		Flags:     flags,
		Exec: func(ctx context.Context, args []string) error {
			res := make(chan error, 1)
			go cmd.handle(args, res)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-res:
				return err
			}
		},
	}

	root.Subcommands = append(root.Subcommands, (*ff.Command)(cmd))
}

func (cmd *createCmd) handle(args []string, res chan<- error) {
	defer close(res)

	home, err := os.UserHomeDir()

	if err != nil {
		res <- err
		return
	}

	labelFlag, _ := cmd.Flags.GetFlag("label")
	hierarchy := strings.Split(labelFlag.GetValue(), ",")
	path := home + "/.anchor" + "/" + strings.Join(hierarchy, ".")

	fh, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		res <- err
		return
	}

	_, err = url.ParseRequestURI(args[1])

	if err != nil {
		res <- ErrInvalidURL
		return
	}

	defer fh.Close()
	_, err = fmt.Fprintf(fh, "%q %q\n", args[0], args[1])

	if err != nil {
		res <- err
	}
}
