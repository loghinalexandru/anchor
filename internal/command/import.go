package command

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/loghinalexandru/anchor/internal/bookmark"
	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/peterbourgon/ff/v4"
	"github.com/virtualtam/netscape-go/v2"
)

var (
	ErrInvalidImportFile = errors.New("invalid import file")
)

type importCmd ff.Command

func newImport(rootFlags *ff.FlagSet) *importCmd {
	var cmd importCmd

	flags := ff.NewFlagSet("import").SetParent(rootFlags)
	cmd = importCmd{
		Name:      "import",
		Usage:     "import",
		ShortHelp: "import bookmarks from a file",
		Flags:     flags,
		Exec:      cmd.handle,
	}

	return &cmd
}

func (*importCmd) handle(_ context.Context, args []string) error {
	dir, err := config.RootDir()
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return ErrInvalidImportFile
	}

	fh, err := os.Open(args[0])
	if err != nil {
		return err
	}

	content, err := io.ReadAll(fh)
	if err != nil {
		return err
	}

	doc, _ := netscape.Unmarshal(content)
	err = traversal(dir, nil, doc.Root)

	if err != nil {
		return err
	}

	return nil
}

func traversal(rootDir string, labels []string, node netscape.Folder) error {
	userDefined, _ := regexp.MatchString("(?i)bookmark|bar", node.Name)

	if len(node.Bookmarks) > 0 && !userDefined {
		labels = append(labels, node.Name)
	}

	file, err := os.OpenFile(filepath.Join(rootDir, FileFrom(labels)), os.O_APPEND|os.O_CREATE|os.O_RDWR, config.StdFileMode)
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
		err = traversal(rootDir, labels, n)
		if err != nil {
			return err
		}
	}

	return nil
}
