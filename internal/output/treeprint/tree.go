package treeprint

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"

	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/xlab/treeprint"
)

const (
	msgMetadata = "%d\u2693"
)

type node struct {
	parent string
	lines  int
}

func Generate(fsys fs.FS) string {
	var hierarchy []map[string]node

	dd, _ := fs.ReadDir(fsys, ".")
	for _, d := range dd {
		if d.IsDir() {
			continue
		}

		var counter int
		fh, err := fsys.Open(d.Name())
		if err == nil {
			counter, _ = lineCounter(fh)
			_ = fh.Close()
		}

		labels := []string{""}
		labels = append(labels, strings.Split(d.Name(), config.StdSeparator)...)
		for i, l := range labels[1:] {
			if len(hierarchy) <= i {
				hierarchy = append(hierarchy, map[string]node{})
			}

			switch i {
			case len(labels) - 2:
				hierarchy[i][l] = node{
					parent: labels[i],
					lines:  counter,
				}
			default:
				if _, ok := hierarchy[i][l]; !ok {
					hierarchy[i][l] = node{
						parent: labels[i],
					}
				}
			}
		}
	}

	return treePrint(hierarchy)
}

func treePrint(hierarchy []map[string]node) string {
	var prev map[string]treeprint.Tree
	var curr map[string]treeprint.Tree

	tree := treeprint.NewWithRoot(filepath.Base(config.RootDir()))
	for _, lvl := range hierarchy {
		curr = make(map[string]treeprint.Tree)
		for _, k := range keys(lvl) {
			if lvl[k].parent == "" {
				br := tree.AddBranch(k)
				curr[k] = br
			} else {
				br := prev[lvl[k].parent]
				curr[k] = br.AddBranch(k)
			}

			if lvl[k].lines > 0 {
				curr[k].SetMetaValue(fmt.Sprintf(msgMetadata, lvl[k].lines))
			}
		}
		prev = curr
	}

	return tree.String()
}

func keys(lvl map[string]node) []string {
	var index int
	keys := make([]string, len(lvl))

	for k := range lvl {
		keys[index] = k
		index++
	}

	slices.Sort(keys)
	return keys
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
