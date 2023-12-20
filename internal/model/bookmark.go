package model

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrDuplicateBookmark = errors.New("duplicate bookmark line")
	ErrInvalidBookmark   = errors.New("cannot parse bookmark: arguments mismatch")
)

type Bookmark struct {
	title  string
	url    string
	client *http.Client
}

func NewBookmark(rawURL string, opts ...func(*Bookmark)) (*Bookmark, error) {
	_, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, err
	}

	res := &Bookmark{
		url:    rawURL,
		client: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(res)
	}

	if res.title == "" {
		res.title = res.fetchTitle()
	}

	return res, nil
}

func WithTitle(title string) func(*Bookmark) {
	return func(b *Bookmark) {
		if title != "" {
			b.title = title
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

func BookmarkLine(line string) (*Bookmark, error) {
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
		return nil, ErrInvalidBookmark
	}

	name, _ := strconv.Unquote(parts[0])
	rawURL, _ := strconv.Unquote(parts[1])

	return &Bookmark{
		title: strings.TrimSpace(name),
		url:   strings.TrimSpace(rawURL),
	}, nil
}

var titleRegexp = regexp.MustCompile(`<title>(?P<title>.+?)</title>`)

// fetchTitle makes a http request to get the html from b.url is no html <title> tag or an error occurs.
// If no html <title> tag is present or an error occurs, returns b.url.
func (b *Bookmark) fetchTitle() string {
	result := b.url

	req, err := http.NewRequest("GET", b.url, nil)
	if err != nil {
		return result
	}

	res, err := b.client.Do(req)
	if err != nil {
		return result
	}
	defer res.Body.Close()

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return result
	}

	match := titleRegexp.FindSubmatch(content)

	if len(match) == 0 {
		return result
	}

	return string(match[1])
}

func (b *Bookmark) String() string {
	return fmt.Sprintf("%q %q\n", b.title, b.url)
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

	exp := regexp.MustCompile(fmt.Sprintf(`(?im)\s.%s.$`, regexp.QuoteMeta(b.url)))
	if exp.Match(content) {
		return fmt.Errorf("%s: %w", b.url, ErrDuplicateBookmark)
	}

	_, err = fmt.Fprint(rw, b.String())
	return err
}

func (b *Bookmark) Update(title string) {
	b.title = title
}

func (b *Bookmark) Title() string {
	return b.title
}

func (b *Bookmark) Description() string {
	return b.url
}

func (b *Bookmark) FilterValue() string {
	return b.title
}
