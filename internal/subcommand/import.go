package subcommand

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/loghinalexandru/anchor/internal/bookmark"
	"github.com/peterbourgon/ff/v4"
	"github.com/virtualtam/netscape-go/v2"
)

var (
	ErrInvalidImportFile = errors.New("invalid import file")
)

type importCmd ff.Command

func RegisterImport(root *ff.Command, rootFlags *ff.FlagSet) {
	var cmd *importCmd

	flags := ff.NewFlagSet("import").SetParent(rootFlags)
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

	rootDir, err := rootDir()
	if err != nil {
		res <- err
		return
	}

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
	err = traversal(rootDir, nil, doc.Root)

	if err != nil {
		res <- err
	}
}

func traversal(rootDir string, labels []string, node netscape.Folder) error {
	userDefined, _ := regexp.MatchString("(?i)bookmark|bar", node.Name)

	if len(node.Bookmarks) > 0 && !userDefined {
		labels = append(labels, node.Name)
	}

	path := filepath.Join(rootDir, formatLabels(labels))
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, fs.ModePerm)
	if err != nil {
		return err
	}

	for _, b := range node.Bookmarks {
		entry, err := bookmark.New(b.Title, b.URL)
		if err != nil {
			return err
		}

		err = entry.Write(file)
		if err != nil && !errors.Is(err, bookmark.ErrDuplicate) {
			return err
		}

		if errors.Is(err, bookmark.ErrDuplicate) {
			fmt.Println(err)
		}
	}

	err = file.Close()
	if err != nil {
		return err
	}

	for _, n := range node.Subfolders {
		err := traversal(rootDir, labels, n)

		if err != nil {
			return err
		}
	}

	return nil
}
