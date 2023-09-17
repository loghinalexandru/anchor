package subcommand

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterbourgon/ff/v4"
	"github.com/xlab/treeprint"
)

const (
	msgMetadata = "%d\u2693"
)

type treeCmd ff.Command

type label struct {
	parent    string
	lineCount int
}

func RegisterTree(root *ff.Command, rootFlags *ff.FlagSet) {
	var cmd *treeCmd

	flags := ff.NewFlagSet("tree").SetParent(rootFlags)
	cmd = &treeCmd{
		Name:      "tree",
		Usage:     "tree",
		ShortHelp: "list available labels in a tree structure",
		Flags:     flags,
		Exec: func(ctx context.Context, args []string) error {
			res := make(chan error, 1)
			go cmd.handle(res)

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

func (*treeCmd) handle(res chan<- error) {
	defer close(res)

	dir, err := rootDir()
	if err != nil {
		res <- err
		return
	}

	dd, err := os.ReadDir(dir)
	if err != nil {
		res <- err
		return
	}

	var hierarchy []map[string]label
	for _, d := range dd {
		if d.IsDir() {
			continue
		}

		fh, err := os.Open(filepath.Join(dir, d.Name()))
		if err != nil {
			res <- err
			return
		}

		c, err := lineCounter(fh)
		err = errors.Join(err, fh.Close())
		if err != nil {
			res <- err
			return
		}

		labels := []string{""}
		labels = append(labels, strings.Split(d.Name(), stdSeparator)...)
		for i, l := range labels[1:] {
			if len(hierarchy) <= i {
				hierarchy = append(hierarchy, map[string]label{})
			}

			switch i {
			case len(labels) - 2:
				hierarchy[i][l] = label{
					parent:    labels[i],
					lineCount: c,
				}
			default:
				if _, ok := hierarchy[i][l]; !ok {
					hierarchy[i][l] = label{
						parent: labels[i],
					}
				}
			}
		}
	}

	var prev map[string]treeprint.Tree
	var curr map[string]treeprint.Tree
	tree := treeprint.NewWithRoot(stdDir)
	for _, lvl := range hierarchy {
		curr = make(map[string]treeprint.Tree)
		for k, v := range lvl {
			if v.parent == "" {
				br := tree.AddBranch(k)
				curr[k] = br
			} else {
				br := prev[v.parent]
				curr[k] = br.AddBranch(k)
			}

			if v.lineCount > 0 {
				curr[k].SetMetaValue(fmt.Sprintf(msgMetadata, v.lineCount))
			}
		}
		prev = curr
	}

	fmt.Print(tree.String())
}
