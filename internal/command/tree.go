package command

import (
	"context"
	"fmt"
	"os"

	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/loghinalexandru/anchor/internal/output/treeprint"
	"github.com/peterbourgon/ff/v4"
)

const (
	treeName = "tree"
)

type treeCmd struct{}

func (tree *treeCmd) manifest(parent *ff.FlagSet) *ff.Command {
	return &ff.Command{
		Name:      treeName,
		Usage:     "anchor tree",
		ShortHelp: "list available labels in a tree structure",
		Flags:     ff.NewFlagSet("tree").SetParent(parent),
		Exec:      tree.handle,
	}
}

func (*treeCmd) handle(context.Context, []string) error {
	dir := config.RootDir()

	dd := os.DirFS(dir)
	fmt.Print(treeprint.Generate(dd))

	return nil
}
