package bookmark

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNew(t *testing.T) {
	t.Parallel()

	title := "test-title"
	url := "https://google.com"

	got, err := New(title, url)
	if err != nil {
		t.Fatalf("unexpected error; got %q", err)
	}

	if got.client == nil {
		t.Error("missing http client")
	}

	if got.Title != title {
		t.Error(cmp.Diff(title, got.Title))
	}

	if got.URL != url {
		t.Error(cmp.Diff(url, got.URL))
	}
}

func TestNewWithInvalidURL(t *testing.T) {
	t.Parallel()

	title := "test-title"
	url := "invalid-url"

	_, err := New(title, url)
	if err == nil {
		t.Error("expected error not found")
	}
}

func TestNewFromLine(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		input string
		title string
		url   string
	}{
		{
			"\r\n\"Outlook\" \"https://outlook.live.com/mail/0/\"\n\r",
			"Outlook",
			"https://outlook.live.com/mail/0/",
		},
		{
			`"Gmail"Test""" "https://accounts.google.com/b/0/AddMailService"   `,
			"GmailTest",
			"https://accounts.google.com/b/0/AddMailService"},
		{
			`"YouTube" "https://youtube.com/"`,
			"YouTube",
			"https://youtube.com/",
		},
	}

	for _, c := range tsc {
		t.Run(c.input, func(t *testing.T) {
			bookmark, err := NewFromLine(c.input)
			if err != nil {
				t.Fatalf("unexpected error; got %q", err)
			}

			if !cmp.Equal(bookmark.Title, c.title) {
				t.Error(cmp.Diff(c.title, bookmark.Title))
			}

			if !cmp.Equal(bookmark.URL, c.url) {
				t.Error(cmp.Diff(c.url, bookmark.URL))
			}
		})
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		title string
		url   string
		want  string
	}{
		{
			title: "Test Title",
			url:   "https://google.com",
			want:  `"Test Title" "https://google.com"`,
		},
		{
			title: `Test "Title" "Test Title Two`,
			url:   "https://google.com",
			want:  `"Test Title Test Title Two" "https://google.com"`,
		},
	}

	for _, c := range tsc {
		t.Run(c.title, func(t *testing.T) {
			bookmark, err := New(c.title, c.url)
			if err != nil {
				t.Fatalf("unexpected error; got %q", err)
			}

			if !cmp.Equal(c.want, bookmark.String()) {
				t.Errorf("wrong serialization: want %s , got: %s", c.want, bookmark.String())
			}
		})
	}
}

func TestTitleFromURL(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		input string
		want  string
	}{
		{
			input: "<meta>stuff></meta><title>Test Title</title><body>test</body>",
			want:  "Test Title",
		},
		{
			input: "<meta>stuff<title>Test Title</title></meta>",
			want:  "Test Title",
		},
		{
			input: `<title>Test - "quote" <title> Title</title>`,
			want:  `Test - "quote" <title> Title`,
		},
		{
			input: "<title>First Test Title</title><title>Second Test Title</title>",
			want:  "First Test Title",
		},
	}

	for _, c := range tcs {
		t.Run(c.input, func(t *testing.T) {
			bookmark, err := New("test", "https://google.com")
			if err != nil {
				t.Fatalf("unexpected error; got %q", err)
			}

			bookmark.client = newTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					Body: io.NopCloser(bytes.NewBufferString(c.input)),
				}
			})

			err = bookmark.TitleFromURL(context.Background())
			if err != nil {
				t.Fatalf("unexpected error; got %q", err)
			}

			if !cmp.Equal(bookmark.Title, c.want) {
				t.Error(cmp.Diff(c.want, bookmark.Title))
			}
		})
	}
}

func TestTitleFromURLWithError(t *testing.T) {
	t.Parallel()

	bookmark, err := New("test-title", "https://google.com")

	if err != nil {
		t.Fatalf("unexpected error; got %q", err)
	}

	bookmark.client = newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			Body: io.NopCloser(bytes.NewBufferString("")),
		}
	})

	err = bookmark.TitleFromURL(context.Background())

	if !errors.Is(err, ErrInvalidTitle) {
		t.Errorf("unexpected error; got %q", err)
	}
}

func TestWrite(t *testing.T) {
	t.Parallel()

	t.TempDir()

	output := filepath.Join(t.TempDir(), t.Name())
	fh, err := os.Create(output)
	if err != nil {
		t.Fatalf("unexpected error; got %q", err)
	}

	want := "\"test-title \\\\n test\" \"https://google.com\"\n"
	bookmark, err := New("test-title \\n test", "https://google.com")
	if err != nil {
		t.Fatalf("unexpected error; got %q", err)
	}

	err = bookmark.Write(fh)
	if err != nil {
		t.Fatalf("unexpected error; got %q", err)
	}

	err = fh.Close()
	if err != nil {
		t.Fatalf("unexpected error; got %q", err)
	}

	content, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("unexpected error; got %q", err)
	}

	got := string(content)

	if !cmp.Equal(got, want) {
		t.Error(cmp.Diff(want, got))
	}
}

type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func newTestClient(fn roundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}
