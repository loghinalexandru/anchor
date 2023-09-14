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

	rep := strings.NewReplacer("\"", "", "\n", "", "\r", "")
	return &Bookmark{
		title: rep.Replace(title),
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

func (b *Bookmark) String() string {
	return fmt.Sprintf("%q %q", b.title, b.url.String())
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

func (b *Bookmark) Write(rw io.ReadWriteSeeker) error {
	_, err := rw.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	content, err := io.ReadAll(rw)
	if err != nil {
		return err
	}

	exp := regexp.MustCompile(fmt.Sprintf(`(?im)\s.%s.$`, regexp.QuoteMeta(b.url.String())))
	if exp.Match(content) {
		return fmt.Errorf("%s: %w", b.url, ErrDuplicate)
	}

	_, err = fmt.Fprintln(rw, b.String())
	return err
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
		fmt.Println(len(res))
		fmt.Println(line)
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
