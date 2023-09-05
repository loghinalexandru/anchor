package subcommand

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/peterbourgon/ff/v4"
)

var (
	ErrInvalidURL   = errors.New("not a valid url")
	ErrDuplicate    = errors.New("duplicate")
	ErrInvalidTitle = errors.New("could not infer title and no flag was set")
)

type createCmd ff.Command

func RegisterCreate(root *ff.Command, rootFlags *ff.CoreFlags) {
	var cmd *createCmd
	flags := ff.NewFlags("create").SetParent(rootFlags)
	_ = flags.String('l', "label", "root", "add label in order of appearance")
	_ = flags.String('t', "title", "", "add custom title")

	cmd = &createCmd{
		Name:      "create",
		Usage:     "crate",
		ShortHelp: "add a bookmark with set labels",
		Flags:     flags,
		Exec: func(ctx context.Context, args []string) error {
			res := make(chan error, 1)
			go cmd.handle(ctx, args, res)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-res:
				return err
			}
		},
	}

	root.Subcommands = append(root.Subcommands, (*ff.Command)(cmd))
}

func (cmd *createCmd) handle(ctx context.Context, args []string, res chan<- error) {
	defer close(res)

	labelFlag, _ := cmd.Flags.GetFlag("label")
	titleFlag, _ := cmd.Flags.GetFlag("title")
	home, err := os.UserHomeDir()

	if err != nil {
		res <- err
		return
	}

	hierarchy := strings.Split(labelFlag.GetValue(), ",")
	path := home + "/.anchor" + "/" + strings.Join(hierarchy, ".")

	fh, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		res <- err
		return
	}

	defer fh.Close()
	url, err := url.ParseRequestURI(args[0])
	if err != nil {
		res <- ErrInvalidURL
		return
	}

	content, _ := io.ReadAll(fh)
	if match, _ := regexp.Match(url.String(), content); match {
		res <- ErrDuplicate
		return
	}

	title, err := title(ctx, args[0], titleFlag.GetValue())
	if err != nil {
		res <- err
		return
	}

	if title == "" {
		res <- ErrInvalidTitle
		return
	}

	_, err = fmt.Fprintf(fh, "%q %q\n", title, args[0])

	if err != nil {
		res <- err
	}
}

func title(ctx context.Context, url string, fallbackValue string) (string, error) {
	titleMatch := regexp.MustCompile("<title>(?P<title>.*)</title>")
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	page, err := io.ReadAll(res.Body)
	if err != nil {
		return "", nil
	}

	m := titleMatch.FindSubmatch(page)
	if len(m) == 0 {
		return fallbackValue, nil
	}

	return string(m[1]), nil
}
