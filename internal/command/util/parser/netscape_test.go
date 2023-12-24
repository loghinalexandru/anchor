package parser_test

import (
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
			name: "gan",
			content: `"Introduction to GANs with Python and TensorFlow" "https://stackabuse.com/introduction-to-gans-with-python-and-tensorflow/"
"sklearn.datasets.fetch_lfw_people — scikit-learn 0.24.1 documentation" "https://scikit-learn.org/stable/modules/generated/sklearn.datasets.fetch_lfw_people.html"
"A Beginner's Guide to Generative Adversarial Networks (GANs) | Pathmind" "https://wiki.pathmind.com/generative-adversarial-network-gan"
"mnist-gan/gan.py at master · gtoubassi/mnist-gan" "https://github.com/gtoubassi/mnist-gan/blob/master/gan.py"
"GitHub - soumith/ganhacks: starter from \"How to Train a GAN?\" at NIPS2016" "https://github.com/soumith/ganhacks"
`,
		},
		{
			name: "gan.research-papers",
			content: `"GAN - 2014 paper" "https://arxiv.org/pdf/1406.2661.pdf"
"LSGAN.pdf" "https://arxiv.org/pdf/1611.04076.pdf"
"Internal Covariate Shift.pdf" "https://arxiv.org/pdf/1502.03167.pdf"
"https://arxiv.org/pdf/1903.06048.pdf" "https://arxiv.org/pdf/1903.06048.pdf"
"https://arxiv.org/pdf/1802.05957.pdf" "https://arxiv.org/pdf/1802.05957.pdf"
`,
		},
	}

	input, err := os.ReadFile("testdata/complex.input")
	if err != nil {
		t.Fatalf("unexpected error when reading input file; got %s", err)
	}

	doc, err := netscape.Unmarshal(input)
	if err != nil {
		t.Fatalf("unexpected error when parsing input file; got %s", err)
	}

	err = parser.Traversal(dir, nil, doc.Root)
	if err != nil {
		t.Fatalf("unexpected error; got %s", err)
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
			t.Errorf(fmt.Sprintf("could not find file name %q", d.Name()))
			return nil
		}

		got, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if want[testCase].content != string(got) {
			t.Errorf("mismatch content; (-want +got):\n %s", cmp.Diff(want[testCase].content, string(got)))
		}

		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error; got %s", err)
	}
}
