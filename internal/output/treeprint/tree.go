package treeprint

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/xlab/treeprint"
)

const (
	msgMetadata = "%d\u2693"
)

func Generate(fsys fs.FS) string {
	known := map[string]treeprint.Tree{}
	tree := treeprint.NewWithRoot(filepath.Base(config.RootDir()))

	dd, _ := fs.ReadDir(fsys, ".")
	for _, d := range dd {
		if d.IsDir() {
			continue
		}

		var lineCount int
		f, err := fsys.Open(d.Name())
		if err == nil {
			lineCount = lineCounter(f)
			f.Close()
		}

		labels := strings.Split(d.Name(), config.StdSeparator)
		if _, ok := known[labels[0]]; !ok {
			known[labels[0]] = branch(tree, lineCount, labels[0], len(labels) == 1)
		}

		for i := 1; i < len(labels); i++ {
			curr := strings.Join(labels[:i+1], config.StdSeparator)
			if _, ok := known[curr]; !ok {
				prev := strings.Join(labels[:i], config.StdSeparator)
				known[curr] = branch(known[prev], lineCount, labels[i], i == len(labels)-1)
			}
		}
	}

	return tree.String()
}

func branch(root treeprint.Tree, lineCount int, label string, leaf bool) treeprint.Tree {
	if leaf {
		return root.AddMetaBranch(fmt.Sprintf(msgMetadata, lineCount), label)
	}

	return root.AddBranch(label)
}

func lineCounter(r io.Reader) int {
	var res int
	buf := make([]byte, 32*1024)
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		res += bytes.Count(buf[:c], lineSep)

		if err != nil {
			return res
		}
	}
}
