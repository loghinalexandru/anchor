package bookmark

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var (
	ErrDuplicate    = errors.New("duplicate bookmark")
	ErrArgsMismatch = errors.New("mismatch in line arguments")
	ErrInvalidTitle = errors.New("could not determine title from URL")
)

type Bookmark struct {
	Name   string
	URL    string
	client *http.Client
}

func New(name string, rawURL string, opts ...func(*Bookmark)) (*Bookmark, error) {
	_, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, err
	}

	result := &Bookmark{
		Name:   sanitize(name),
		URL:    rawURL,
		client: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(result)
	}

	return result, nil
}

func WithClient(client *http.Client) func(*Bookmark) {
	return func(b *Bookmark) {
		b.client = client
	}
}

func NewFromLine(line string) (*Bookmark, error) {
	line = strings.Trim(line, " \"\r\n")
	parts := strings.Split(line, `" "`)

	if len(parts) != 2 {
		return nil, ErrArgsMismatch
	}

	return New(parts[0], parts[1])
}

func (b *Bookmark) String() string {
	return fmt.Sprintf("%q %q", b.Name, b.URL)
}

func (b *Bookmark) TitleFromURL(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", b.URL, nil)
	if err != nil {
		return err
	}

	res, err := b.client.Do(req)
	if err != nil {
		return err
	}

	page, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = res.Body.Close()
	if err != nil {
		return err
	}

	title := findTitle(page)
	if title == "" {
		return ErrInvalidTitle
	}

	b.Name = title
	return nil
}

func (b *Bookmark) Write(rw io.ReadWriteSeeker) error {
	_, err := rw.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	content, err := io.ReadAll(rw)
	if err != nil {
		return err
	}

	exp := regexp.MustCompile(fmt.Sprintf(`(?im)\s.%s.$`, regexp.QuoteMeta(b.URL)))
	if exp.Match(content) {
		return fmt.Errorf("%s: %w", b.URL, ErrDuplicate)
	}

	_, err = fmt.Fprintln(rw, b.String())
	return err
}

func (b *Bookmark) Title() string {
	return b.Name
}

func (b *Bookmark) Description() string {
	return b.URL
}

func (b *Bookmark) FilterValue() string {
	return b.Name
}

func sanitize(input string) string {
	repl := strings.NewReplacer("\n", "", "\r", "", "\"", "")
	return repl.Replace(input)
}

var titleRegexp = regexp.MustCompile(`<title>(?P<title>.+?)</title>`)

func findTitle(content []byte) string {
	match := titleRegexp.FindSubmatch(content)

	if len(match) == 0 {
		return ""
	}

	return string(match[1])
}
