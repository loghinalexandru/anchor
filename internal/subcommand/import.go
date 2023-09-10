package subcommand

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/loghinalexandru/anchor/internal/bookmark"
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

	dir, _ := c.Flags.GetFlag("root-dir")
	home, err := os.UserHomeDir()
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

	path := filepath.Join(home, dir.GetValue())
	doc, _ := netscape.Unmarshal(content)
	err = traversal(path, "", doc.Root)

	if err != nil {
		res <- err
	}
}

// Refactor this
func traversal(basePath string, fileName string, node netscape.Folder) error {
	isRoot, _ := regexp.MatchString("(?i)bookmark|bar", node.Name)

	if len(node.Bookmarks) > 0 && !isRoot {
		label := strings.ReplaceAll(node.Name, " ", "")
		label = strings.ToLower(label)

		if fileName != "" {
			fileName = fmt.Sprintf("%s.%s", fileName, label)
		} else {
			fileName = label
		}
	}

	for _, b := range node.Bookmarks {
		var filePath string
		entry, err := bookmark.New(b.Title, b.URL)
		if err != nil {
			return err
		}

		if fileName == "" {
			filePath = filepath.Join(basePath, "root")
		} else {
			filePath = filepath.Join(basePath, fileName)
		}

		_, err = bookmark.Append(*entry, filePath)
		if err != nil && !errors.Is(err, bookmark.ErrDuplicate) {
			return err
		}
	}

	for _, n := range node.Subfolders {
		err := traversal(basePath, fileName, n)

		if err != nil {
			return err
		}
	}

	return nil
}
