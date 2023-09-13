package bookmark

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

const (
	regexpEndOfLine = `(?im)\s.%s.$`
)

var (
	ErrDuplicate    = errors.New("duplicate bookmark")
	ErrArgsMismatch = errors.New("mismatch in line arguments")
	ErrInvalidTitle = errors.New("could not infer title and no flag was set")
)

type Bookmark struct {
	title string
	url   *url.URL
}

func New(title string, rawURL string) (*Bookmark, error) {
	url, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, err
	}

	return &Bookmark{
		title: title,
		url:   url,
	}, nil
}

func NewFromLine(line string) (*Bookmark, error) {
	title, url, err := Parse(line)
	if err != nil {
		return nil, err
	}

	return New(title, url)
}

func (b Bookmark) String() string {
	return fmt.Sprintf("%q %q", b.title, strings.Trim(b.url.String(), " \r\n"))
}

func (b *Bookmark) TitleFromURL(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", b.url.String(), nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	page, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	res.Body.Close()

	title := findTitle(page)
	if title == "" {
		return ErrInvalidTitle
	}

	b.title = title
	return nil
}

func Append(b Bookmark, filePath string) (int, error) {
	fh, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, fs.ModePerm)
	if err != nil {
		return 0, err
	}

	content, _ := io.ReadAll(fh)
	defer fh.Close()

	if ok := matchEndOfLines(content, b.url.String()); ok {
		return 0, fmt.Errorf("%q: %w", b.url.String(), ErrDuplicate)
	}

	return fmt.Fprintln(fh, b.String())
}

func Parse(line string) (title, url string, err error) {
	quoted := false
	line = strings.Trim(line, " \r\n")

	res := strings.FieldsFunc(line, func(r rune) bool {
		if r == '"' {
			quoted = !quoted
		}
		return !quoted && r == ' '
	})

	if len(res) != 2 {
		return title, url, ErrArgsMismatch
	}

	title = strings.Trim(res[0], " \"\r\n")
	url = strings.Trim(res[1], " \"\r\n")

	return title, url, nil
}

func findTitle(content []byte) string {
	titleMatch := regexp.MustCompile("<title>(?P<title>.*)</title>")
	match := titleMatch.FindSubmatch(content)

	if len(match) == 0 {
		return ""
	}

	return string(match[1])
}

func matchEndOfLines(content []byte, pattern string) bool {
	regex := regexp.MustCompile(fmt.Sprintf(regexpEndOfLine, regexp.QuoteMeta(pattern)))
	return regex.Match(content)
}
