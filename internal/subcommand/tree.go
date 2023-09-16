package subcommand

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/peterbourgon/ff/v4"
	"github.com/xlab/treeprint"
)

type treeCmd ff.Command

func RegisterTree(root *ff.Command, rootFlags *ff.FlagSet) {
	var cmd *treeCmd

	flags := ff.NewFlagSet("tree").SetParent(rootFlags)
	cmd = &treeCmd{
		Name:      "tree",
		Usage:     "tree",
		ShortHelp: "list available labels in tree structure",
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

	var hierarchy []map[string]string
	for _, d := range dd {
		if d.IsDir() {
			continue
		}

		labels := strings.Split(d.Name(), ".")
		for i, l := range labels {
			if len(hierarchy) <= i {
				hierarchy = append(hierarchy, make(map[string]string))
			}

			if i == 0 {
				hierarchy[i][l] = ""
			} else {
				hierarchy[i][l] = labels[i-1]
			}
		}
	}

	var prev map[string]treeprint.Tree
	var curr map[string]treeprint.Tree
	tree := treeprint.NewWithRoot(defaultDir)
	for _, lvl := range hierarchy {
		curr = make(map[string]treeprint.Tree)
		for k, v := range lvl {
			if v == "" {
				br := tree.AddBranch(k)
				curr[k] = br
			} else {
				br := prev[v]
				curr[k] = br.AddBranch(k)
			}
		}
		prev = curr
	}

	fmt.Println(tree.String())
}

// Add line counter for each label
// func lineCounter(r io.Reader) (int, error) {
// 	var res int

// 	buf := make([]byte, 32*1024)
// 	lineSep := []byte{'\n'}

// 	for {
// 		c, err := r.Read(buf)
// 		res += bytes.Count(buf[:c], lineSep)

// 		switch {
// 		case err == io.EOF:
// 			return res, nil
// 		case err != nil:
// 			return res, err
// 		}
// 	}
// }
