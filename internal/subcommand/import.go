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

	rootDir := filepath.Join(home, dir.GetValue())
	doc, _ := netscape.Unmarshal(content)
	err = traversal(rootDir, "", doc.Root)

	if err != nil {
		res <- err
	}
}

func traversal(rootDir string, labels string, node netscape.Folder) error {
	isRoot, _ := regexp.MatchString("(?i)bookmark|bar", node.Name)

	if len(node.Bookmarks) > 0 && !isRoot {
		if labels == "" {
			labels = node.Name
		} else {
			labels = fmt.Sprintf("%s.%s", labels, node.Name)
		}
	}

	path := filepath.Join(rootDir, format(labels))

	for _, b := range node.Bookmarks {
		entry, err := bookmark.New(b.Title, b.URL)
		if err != nil {
			return err
		}

		_, err = bookmark.Append(*entry, path)
		if err != nil && !errors.Is(err, bookmark.ErrDuplicate) {
			return err
		}
	}

	for _, n := range node.Subfolders {
		err := traversal(rootDir, labels, n)

		if err != nil {
			return err
		}
	}

	return nil
}

// Fix this random root
func format(labels string) string {
	if labels == "" {
		return "root"
	}

	exp := regexp.MustCompile(`[^a-z0-9-\.]`)
	return exp.ReplaceAllString(strings.ToLower(labels), "")
}
