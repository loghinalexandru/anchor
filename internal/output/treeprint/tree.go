package treeprint

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/xlab/treeprint"
)

const (
	msgMetadata = "%d\u2693"
)

func Generate(fsys fs.FS) string {
	known := map[string]treeprint.Tree{}
	tree := treeprint.NewWithRoot(config.StdDirName)

	// For each dir in fsys read and compute number of lines together
	// with the tree from the flat structure of the file names.
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

		// Split the file name by config.StdLabelSeparator and
		// create a tree root node from first label.
		labels := strings.Split(d.Name(), config.StdLabelSeparator)
		if _, ok := known[labels[0]]; !ok {
			known[labels[0]] = branch(tree, lineCount, labels[0], len(labels) == 1)
		}

		// Go over the rest of the labels and create a new tree node
		// with the previous labels as parent if it was not seen before.
		for i := 1; i < len(labels); i++ {
			curr := strings.Join(labels[:i+1], config.StdLabelSeparator)
			if _, ok := known[curr]; !ok {
				prev := strings.Join(labels[:i], config.StdLabelSeparator)
				known[curr] = branch(known[prev], lineCount, labels[i], i == len(labels)-1)
			}
		}
	}

	return tree.String()
}

// Add line count metadata if the label is the last one.
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
