package subcommand

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/peterbourgon/ff/v4"
	"github.com/virtualtam/netscape-go/v2"
)

var (
	ErrInvalidImportFile = errors.New("invalid import file")
)

type importCmd ff.Command

func RegisterImport(root *ff.Command, rootFlags *ff.CoreFlags) {
	var cmd *importCmd
	flags := ff.NewFlags("import").SetParent(rootFlags)

	cmd = &importCmd{
		Name:      "import",
		Usage:     "import",
		ShortHelp: "import bookmarks from a file",
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

func (c *importCmd) handle(args []string, res chan<- error) {
	defer close(res)

	if len(args) == 0 {
		res <- ErrInvalidImportFile
		return
	}

	fh, err := os.Open(args[0])
	if err != nil {
		res <- err
		return
	}

	content, err := io.ReadAll(fh)
	if err != nil {
		res <- err
		return
	}

	doc, _ := netscape.Unmarshal(content)
	fmt.Print(doc)
}
