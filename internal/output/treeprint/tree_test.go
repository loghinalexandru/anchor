package treeprint_test

import (
	"testing"
	"testing/fstest"

	"github.com/loghinalexandru/anchor/internal/output/treeprint"
)

func TestGenerateNoMetadata(t *testing.T) {
	t.Parallel()

	fs := fstest.MapFS{
		"root":          {},
		"prometheus.go": {},
		"prometheus":    {},
	}

	got := treeprint.Generate(fs)
	t.Error(got)
}
