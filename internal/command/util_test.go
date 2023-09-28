package command

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestFindLines(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		pattern    string
		matchCount int
	}{
		{
			pattern:    "y",
			matchCount: 1,
		},
		{
			pattern:    "ou",
			matchCount: 2,
		},
		{
			pattern:    "aws",
			matchCount: 1,
		},
		{
			pattern:    "o",
			matchCount: 4,
		},
		{
			pattern:    "Gmail",
			matchCount: 1,
		},
	}

	file, err := os.Open("testdata/root.txt")
	if err != nil {
		t.Fatalf("unexpected error; got %q", err)
	}

	defer file.Close()

	in, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("unexpected error; got %q", err)

	}

	for _, c := range tsc {
		t.Run(c.pattern, func(t *testing.T) {
			got := FindLines(in, c.pattern)

			if len(got) != c.matchCount {
				t.Fatalf("unexpected match count; want %q, got %q", c.matchCount, len(got))
			}

			for _, match := range got {
				if !strings.Contains(string(match), c.pattern) {
					t.Log(string(match))
					t.Errorf("missing pattern %q in match", c.pattern)
				}
			}
		})
	}
}
