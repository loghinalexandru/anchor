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
	link   string
	client *http.Client
}

func NewBookmark(rawURL string, opts ...func(*Bookmark)) (*Bookmark, error) {
	_, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, err
	}

	res := &Bookmark{
		link:   rawURL,
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
	link, _ := strconv.Unquote(parts[1])

	return &Bookmark{
		title: strings.TrimSpace(name),
		link:  strings.TrimSpace(link),
	}, nil
}

var titleRegexp = regexp.MustCompile(`<title>(?P<title>.+?)</title>`)

// If there is no html <title> tag or an error occurs
// returns the bookmark link.
func (b *Bookmark) fetchTitle() string {
	result := b.link

	req, err := http.NewRequest("GET", b.link, nil)
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
	return fmt.Sprintf("%q %q\n", b.title, b.link)
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

	exp := regexp.MustCompile(fmt.Sprintf(`(?im)\s.%s.$`, regexp.QuoteMeta(b.link)))
	if exp.Match(content) {
		return fmt.Errorf("%s: %w", b.link, ErrDuplicateBookmark)
	}

	_, err = fmt.Fprint(rw, b.String())
	return err
}

func (b *Bookmark) Title() string {
	return b.title
}

func (b *Bookmark) Description() string {
	return b.link
}

func (b *Bookmark) FilterValue() string {
	return b.title
}
