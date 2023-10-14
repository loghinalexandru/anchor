package output

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/xlab/treeprint"
)

const (
	msgMetadata = "%d\u2693"
)

type label struct {
	parent string
	lines  int
}

func Tree(basePath string, dd []os.DirEntry) string {
	var hierarchy []map[string]label
	for _, d := range dd {
		if d.IsDir() {
			continue
		}

		var counter int
		fh, err := os.Open(filepath.Join(basePath, d.Name()))
		if err == nil {
			counter, _ = lineCounter(fh)
			_ = fh.Close()
		}

		labels := []string{""}
		labels = append(labels, strings.Split(d.Name(), config.StdSeparator)...)
		for i, l := range labels[1:] {
			if len(hierarchy) <= i {
				hierarchy = append(hierarchy, map[string]label{})
			}

			switch i {
			case len(labels) - 2:
				hierarchy[i][l] = label{
					parent: labels[i],
					lines:  counter,
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

	return treePrint(hierarchy)
}

func treePrint(hierarchy []map[string]label) string {
	var prev map[string]treeprint.Tree
	var curr map[string]treeprint.Tree
	tree := treeprint.NewWithRoot(filepath.Base(config.RootDir()))
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

			if v.lines > 0 {
				curr[k].SetMetaValue(fmt.Sprintf(msgMetadata, v.lines))
			}
		}
		prev = curr
	}

	return tree.String()
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
