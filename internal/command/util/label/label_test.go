package label

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestName(t *testing.T) {
	t.Parallel()

	tsc := map[string]struct {
		labels []string
		want   string
	}{
		"empty-label": {
			labels: []string{},
			want:   filepath.Join("test", "root"),
		},
		"single-label": {
			labels: []string{"dotnet"},
			want:   filepath.Join("test", "dotnet"),
		},
		"multi-label": {
			labels: []string{"prometheus", "golang"},
			want:   filepath.Join("test", "prometheus.golang"),
		},
	}

	for k, c := range tsc {
		t.Run(k, func(t *testing.T) {
			got := name("test", c.labels)
			if c.want != got {
				t.Errorf("unexpected filename from label(s); want %q, got %q", c.want, got)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	t.Parallel()

	tsc := map[string]struct {
		labels []string
		want   []string
	}{
		"whitespace-label": {
			labels: []string{"  \n"},
			want:   []string{""},
		},
		"single-label": {
			labels: []string{".Dotnet-TEST!@#"},
			want:   []string{"dotnet-test"},
		},
		"multi-label": {
			labels: []string{"!@#$%^&*();'proMethe%us", ".. gola ng"},
			want:   []string{"prometheus", "golang"},
		},
	}

	for k, c := range tsc {
		t.Run(k, func(t *testing.T) {
			got := Format(c.labels)
			if !cmp.Equal(c.want, got) {
				t.Errorf("unexpected format for labels; want %q, got %q", c.want, got)
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
