package subcommand

import (
	"context"
	"os"
	"path/filepath"

	"github.com/loghinalexandru/anchor/internal/bookmark"
	"github.com/peterbourgon/ff/v4"
)

type create struct {
	command *ff.Command
	labels  *[]string
}

func RegisterCreate(root *ff.Command, rootFlags *ff.CoreFlags) {
	var cr create
	var labels []string

	flags := ff.NewFlags("create").SetParent(rootFlags)
	_ = flags.StringSetVar(&labels, 'l', "label", "add labels in order of appearance")
	_ = flags.String('t', "title", "", "add custom title")

	cr = create{
		command: &ff.Command{
			Name:      "create",
			Usage:     "crate",
			ShortHelp: "add a bookmark with set labels",
			Flags:     flags,
			Exec: func(ctx context.Context, args []string) error {
				res := make(chan error, 1)
				go cr.handle(ctx, args, res)

				select {
				case <-ctx.Done():
					return ctx.Err()
				case err := <-res:
					return err
				}
			},
		},
		labels: &labels,
	}

	root.Subcommands = append(root.Subcommands, cr.command)
}

func (c create) handle(ctx context.Context, args []string, res chan<- error) {
	defer close(res)

	titleFlag, _ := c.command.Flags.GetFlag("title")
	dir, _ := c.command.Flags.GetFlag("root-dir")
	home, err := os.UserHomeDir()

	if err != nil {
		res <- err
		return
	}

	b, err := bookmark.New(titleFlag.GetValue(), args[0])
	if err != nil {
		res <- err
		return
	}

	if titleFlag.GetValue() == titleFlag.GetDefault() {
		err = b.TitleFromURL(ctx)

		if err != nil {
			res <- err
			return
		}
	}

	err = validate(*c.labels)
	if err != nil {
		res <- err
		return
	}

	tree := formatLabels(*c.labels)
	path := filepath.Join(home, dir.GetValue(), tree)
	_, err = bookmark.Append(*b, path)

	if err != nil {
		res <- err
	}
}
