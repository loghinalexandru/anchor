package label

import (
	"path/filepath"
	"testing"
)

func TestFilename(t *testing.T) {
	t.Parallel()

	tsc := map[string]struct {
		labels []string
		want   string
	}{
		"empty-label": {
			labels: []string{},
			want:   "root",
		},
		"single-label": {
			labels: []string{"dotnet"},
			want:   "dotnet",
		},
		"multi-label": {
			labels: []string{"prometheus", "golang"},
			want:   "prometheus.golang",
		},
	}

	for k, c := range tsc {
		t.Run(k, func(t *testing.T) {
			got := Filepath(c.labels)
			if c.want != filepath.Base(got) {
				t.Errorf("unexpected filename from labels; want %q, got %q", c.want, filepath.Base(got))
			}
		})
	}
}

func TestValidate(t *testing.T) {
	t.Parallel()

}
