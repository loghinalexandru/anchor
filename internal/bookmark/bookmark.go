package bookmark

import (
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
)

type Bookmark struct {
	Name   string
	URL    string
	client *http.Client
}

func New(rawURL string, opts ...func(*Bookmark)) (*Bookmark, error) {
	_, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, err
	}

	res := &Bookmark{
		URL:    rawURL,
		client: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(res)
	}

	if res.Name == "" {
		res.Name = res.fetchTitle()
	}

	return res, nil
}

func WithTitle(title string) func(*Bookmark) {
	return func(b *Bookmark) {
		if title != "" {
			b.Name = title
		}
	}
}

func WithClient(client *http.Client) func(*Bookmark) {
	return func(b *Bookmark) {
		if client != nil {
			b.client = client
		}
	}
}

func NewFromLine(line string) (*Bookmark, error) {
	var quoted bool
	var prev rune

	line = strings.Trim(line, " \r\n")
	parts := strings.FieldsFunc(line, func(curr rune) bool {
		if curr == '"' && prev != '\\' {
			quoted = !quoted
		}

		prev = curr
		return !quoted && curr == ' '
	})

	if len(parts) != 2 {
		return nil, ErrArgsMismatch
	}

	return &Bookmark{
		Name: strings.Replace(strings.Trim(parts[0], " \""), "\\", "", -1),
		URL:  strings.Replace(strings.Trim(parts[1], " \""), "\\", "", -1),
	}, nil
}

var titleRegexp = regexp.MustCompile(`<title>(?P<title>.+?)</title>`)

func (b *Bookmark) fetchTitle() string {
	result := b.URL

	req, err := http.NewRequest("GET", b.URL, nil)
	if err != nil {
		return result
	}

	res, err := b.client.Do(req)
	if err != nil {
		return result
	}
	defer res.Body.Close()

	page, err := io.ReadAll(res.Body)
	if err != nil {
		return result
	}

	match := titleRegexp.FindSubmatch(page)

	if len(match) == 0 {
		return b.URL
	}

	return string(match[1])
}

func (b *Bookmark) String() string {
	return fmt.Sprintf("%q %q\n", b.Name, b.URL)
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

	_, err = fmt.Fprint(rw, b.String())
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
