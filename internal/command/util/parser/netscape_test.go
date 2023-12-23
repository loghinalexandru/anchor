package parser_test

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/loghinalexandru/anchor/internal/command/util/parser"
	"github.com/virtualtam/netscape-go/v2"
)

type TraversalTest struct {
	name    string
	content string
}

func TestTraversalComplex(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	want := []TraversalTest{
		{
			name:    "gan",
			content: "asd",
		},
	}

	input, err := os.ReadFile("testdata/complex.input")
	if err != nil {
		t.Fatalf("unexpected error when reading input file; got %q", err)
	}

	doc, err := netscape.Unmarshal(input)
	if err != nil {
		t.Fatalf("unexpected error when parsing input file; got %q", err)
	}

	err = parser.Traversal(dir, nil, doc.Root)
	if err != nil {
		t.Fatalf("unexpected error; got %q", err)
	}

	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		testCase := slices.IndexFunc(want, func(tt TraversalTest) bool {
			if tt.name == d.Name() {
				return true
			}

			return false
		})

		if testCase == -1 {
			return errors.New(fmt.Sprintf("could not find file name %q", d.Name()))
		}

		got, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if want[testCase].content != string(got) {
			return errors.New(fmt.Sprintf("unexpected content for file %q; (-want +got):\\n%s", d.Name(), cmp.Diff(want[testCase].content, string(got))))
		}

		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error; got %q", err)
	}
}
