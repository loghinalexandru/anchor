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

func (c *treeCmd) handle(res chan<- error) {
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

		labels := strings.Split(d.Name(), stdSeparator)
		for i, l := range labels {
			if len(hierarchy) <= i {
				hierarchy = append(hierarchy, make(map[string]label))
			}

			switch i {
			case 0:
				hierarchy[i][l] = label{
					lineCount: c,
				}
			case len(labels) - 1:
				hierarchy[i][l] = label{
					parent:    labels[i-1],
					lineCount: c,
				}
			default:
				hierarchy[i][l] = label{
					parent: labels[i-1],
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
				br := tree.AddMetaBranch(fmt.Sprintf(msgMetadata, v.lineCount), k)
				curr[k] = br
			} else {
				br := prev[v.parent]
				if v.lineCount > 0 {
					curr[k] = br.AddMetaBranch(fmt.Sprintf(msgMetadata, v.lineCount), k)
				} else {
					curr[k] = br.AddBranch(k)
				}
			}
		}
		prev = curr
	}

	fmt.Print(tree.String())
}
