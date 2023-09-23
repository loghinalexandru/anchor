package command

import (
	"context"
	"fmt"
	"os"

	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/loghinalexandru/anchor/internal/output"
	"github.com/peterbourgon/ff/v4"
)

type treeCmd ff.Command

func newTree(rootFlags *ff.FlagSet) *treeCmd {
	var cmd *treeCmd

	flags := ff.NewFlagSet("tree").SetParent(rootFlags)
	cmd = &treeCmd{
		Name:      "tree",
		Usage:     "tree",
		ShortHelp: "list available labels in a tree structure",
		Flags:     flags,
		Exec:      handlerMiddleware(cmd.handle),
	}

	return cmd
}

func (*treeCmd) handle(context.Context, []string) error {

	dir, err := config.RootDir()
	if err != nil {
		return err
	}

	dd, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	fmt.Print(output.Tree(dir, dd))
	return nil
}