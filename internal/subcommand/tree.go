package subcommand

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

	tree := treeprint.New()
	for _, d := range dd {
		if d.IsDir() {
			continue
		}

		labels := strings.Split(d.Name(), ".")
		br := tree.Branch()

		for _, l := range labels {
			br = br.AddBranch(l)
		}

		fh, err := os.Open(filepath.Join(dir, d.Name()))
		if err != nil {
			res <- err
			return
		}

		count, err := lineCounter(fh)
		fh.Close()

		if err != nil {
			res <- err
			return
		}

		br.SetMetaValue(count)
	}

	fmt.Println(tree.String())
}

func lineCounter(r io.Reader) (int, error) {
	var res int

	buf := make([]byte, 32*1024)
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		res += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return res, nil
		case err != nil:
			return res, err
		}
	}
}
