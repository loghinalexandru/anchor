package treeprint_test

import (
	"testing"
	"testing/fstest"

	"github.com/google/go-cmp/cmp"
	"github.com/loghinalexandru/anchor/internal/output/treeprint"
)

func TestGenerateNoMetadata(t *testing.T) {
	t.Parallel()

	want := ".anchor\n├── prometheus\n│   ├── algo\n│   └── go\n└── root\n"
	fs := fstest.MapFS{
		"root":            {},
		"prometheus.go":   {},
		"prometheus.algo": {},
		"prometheus":      {},
	}

	got := treeprint.Generate(fs)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("output not matching; (-got, +want):\n %s", diff)
	}
}
