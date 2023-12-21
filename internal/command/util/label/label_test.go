package label

import (
	"errors"
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
			got := Filename(c.labels)
			if c.want != filepath.Base(got) {
				t.Errorf("unexpected filename from label(s); want %q, got %q", c.want, filepath.Base(got))
			}
		})
	}
}

func TestValidateGoodLabels(t *testing.T) {
	t.Parallel()

	tsc := map[string][]string{
		"empty-label":  {},
		"single-label": {"golang"},
		"multi-label":  {"golang", "prometheus", "1-config"},
	}

	for k, c := range tsc {
		t.Run(k, func(t *testing.T) {
			err := validate(c)
			if err != nil {
				t.Errorf("unexpected error on valid label(s); %q", err)
			}
		})
	}
}

func TestValidateBadLabels(t *testing.T) {
	t.Parallel()

	tsc := map[string][]string{
		"empty-label":      {""},
		"whitespace-label": {" \n"},
		"single-label":     {".dotnet"},
		"multi-label":      {"Golang", "#!?prometheus", "1-config"},
	}

	for k, c := range tsc {
		t.Run(k, func(t *testing.T) {
			err := validate(c)
			if !errors.Is(err, ErrInvalidLabel) {
				t.Errorf("missing error on invalid label(s); %q", c)
			}
		})
	}
}
