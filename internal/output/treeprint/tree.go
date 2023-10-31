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

		var lines int
		f, err := fsys.Open(d.Name())
		if err == nil {
			lines = lineCounter(f)
			f.Close()
		}

		labels := strings.Split(d.Name(), config.StdSeparator)
		if _, ok := known[labels[0]]; !ok {
			known[labels[0]] = tree.AddMetaBranch(fmt.Sprintf(msgMetadata, lines), labels[0])
		}

		for i := 1; i < len(labels); i++ {
			curr := strings.Join(labels[:i+1], config.StdSeparator)
			if _, ok := known[curr]; !ok {
				prev := strings.Join(labels[:i], config.StdSeparator)
				if i == len(labels)-1 {
					known[curr] = known[prev].AddMetaBranch(fmt.Sprintf(msgMetadata, lines), labels[i])
				} else {
					known[curr] = known[prev].AddBranch(labels[i])
				}
			}
		}
	}

	return tree.String()
}

func lineCounter(r io.Reader) int {
	var res int
	buf := make([]byte, 32*1024)
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		res += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return res
		case err != nil:
			return res
		}
	}
}
