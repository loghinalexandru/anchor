package bookmark_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/loghinalexandru/anchor/internal/bookmark"
)

func TestNew(t *testing.T) {
	t.Parallel()

	title := "test-title"
	url := "https://google.com"

	got, err := bookmark.New(url, bookmark.WithTitle(title))
	if err != nil {
		t.Fatalf("unexpected error; got %q", err)
	}

	if got.Name != title {
		t.Error(cmp.Diff(title, got.Title))
	}

	if got.URL != url {
		t.Error(cmp.Diff(url, got.URL))
	}
}

func TestNewWithInvalidURL(t *testing.T) {
	t.Parallel()

	url := "invalid-url"

	_, err := bookmark.New(url)
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
			"\r\n\"Outlook\"    \"https://outlook.live.com/mail/0/\"\n\r",
			"Outlook",
			"https://outlook.live.com/mail/0/",
		},
		{
			"\"GmailTest \" \" \"https://accounts.google.com/b/0/AddMailService\"   ",
			"GmailTest",
			"https://accounts.google.com/b/0/AddMailService",
		},
		{
			"\"YouTube \" \"   \"https://youtube.com/\"",
			"YouTube",
			"https://youtube.com/",
		},
	}

	for _, c := range tsc {
		t.Run(c.input, func(t *testing.T) {
			bk, err := bookmark.NewFromLine(c.input)
			if err != nil {
				t.Fatalf("unexpected error; got %q", err)
			}

			if bk.Name != c.title {
				t.Errorf("wrong deserialization: want %s , got: %s", c.title, bk.Name)
			}

			if bk.URL != c.url {
				t.Errorf("wrong deserialization: want %s , got: %s", c.url, bk.URL)
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
			want:  "\"Test Title\" \"https://google.com\"\n",
		},
		{
			title: `Test "Title" "Test Title Two`,
			url:   "https://google.com",
			want:  "\"Test \\\"Title\\\" \\\"Test Title Two\" \"https://google.com\"\n",
		},
	}

	for _, c := range tsc {
		t.Run(c.title, func(t *testing.T) {
			bk, err := bookmark.New(c.url, bookmark.WithTitle(c.title))
			if err != nil {
				t.Fatalf("unexpected error; got %q", err)
			}

			if c.want != bk.String() {
				t.Errorf("wrong serialization: want %s , got: %s", c.want, bk.String())
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
			bk, err := bookmark.New("https://google.com", bookmark.WithClient(newTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					Body: io.NopCloser(bytes.NewBufferString(c.input)),
				}
			})))

			if err != nil {
				t.Fatalf("unexpected error; got %q", err)
			}

			if err != nil {
				t.Fatalf("unexpected error; got %q", err)
			}

			if !cmp.Equal(bk.Name, c.want) {
				t.Error(cmp.Diff(c.want, bk.Name))
			}
		})
	}
}

func TestWrite(t *testing.T) {
	t.Parallel()

	output := filepath.Join(t.TempDir(), t.Name())
	fh, err := os.Create(output)
	if err != nil {
		t.Fatalf("unexpected error; got %q", err)
	}

	title := "test-title \\n \"test\" asd"
	want := "\"test-title \\\\n \\\"test\\\" asd\" \"https://google.com\"\n"
	bk, err := bookmark.New("https://google.com", bookmark.WithTitle(title))
	if err != nil {
		t.Fatalf("unexpected error; got %q", err)
	}

	err = bk.Write(fh)
	if err != nil {
		t.Fatalf("unexpected error; got %q", err)
	}

	err = fh.Close()
	if err != nil {
		t.Fatalf("unexpected error; got %q", err)
	}

	got, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("unexpected error; got %q", err)
	}

	fmt.Printf(string(got))

	if string(got) != want {
		t.Errorf("wrong serialization: want %s , got: %s", want, got)
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
