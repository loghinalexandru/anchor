package command

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/loghinalexandru/anchor/internal/command/util/parser"
	"github.com/loghinalexandru/anchor/internal/config"
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
		Usage:     "anchor import [PATH]",
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
	err = parser.TraverseNode(config.RootDir(), nil, doc.Root)

	if err != nil {
		return err
	}

	return nil
}
