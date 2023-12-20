package command

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/loghinalexandru/anchor/internal/command/util/label"
	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/loghinalexandru/anchor/internal/model"
	"github.com/peterbourgon/ff/v4"
	"github.com/virtualtam/netscape-go/v2"
)

const (
	importName = "import"
)

var (
	ErrInvalidImportFile = errors.New("invalid import file")
)

type importCmd struct{}

func (imp *importCmd) manifest(parent *ff.FlagSet) *ff.Command {
	return &ff.Command{
		Name:      importName,
		Usage:     "import",
		ShortHelp: "import bookmarks from a file",
		Flags:     ff.NewFlagSet("import").SetParent(parent),
		Exec:      imp.handle,
	}
}

func (*importCmd) handle(_ context.Context, args []string) error {
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
	err = traversal(config.RootDir(), nil, doc.Root)

	if err != nil {
		return err
	}

	return nil
}

var toolbarRegexp = regexp.MustCompile("(?i)bookmark|bar")

func traversal(rootDir string, labels []string, node netscape.Folder) error {
	userDefined := toolbarRegexp.MatchString(node.Name)

	if len(node.Bookmarks) > 0 && !userDefined {
		labels = append(labels, node.Name)
	}

	file, err := os.OpenFile(label.Filepath(labels), os.O_APPEND|os.O_CREATE|os.O_RDWR, config.StdFileMode)
	if err != nil {
		return err
	}

	for _, b := range node.Bookmarks {
		entry, err := model.NewBookmark(b.URL, model.WithTitle(b.Title))
		if err != nil {
			return err
		}

		err = entry.Write(file)
		if err != nil && !errors.Is(err, model.ErrDuplicateBookmark) {
			return err
		}

		if errors.Is(err, model.ErrDuplicateBookmark) {
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
