package subcommand

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/loghinalexandru/anchor/internal/regex"
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
	labelFlag, _ := cmd.Flags.GetFlag("label")
	titleFlag, _ := cmd.Flags.GetFlag("title")
	dir, _ := cmd.Flags.GetFlag("root-dir")
	home, err := os.UserHomeDir()

	if err != nil {
		res <- err
		return
	}

	hierarchy := strings.Split(labelFlag.GetValue(), ",")
	path := filepath.Join(home, dir.GetValue(), strings.Join(hierarchy, "."))

	fh, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, fs.ModePerm)
	if err != nil {
		res <- err
		return
	}

	if len(args) == 0 {
		res <- ErrInvalidURL
		return
	}

	bookmark := args[0]
	_, err = url.ParseRequestURI(bookmark)
	if err != nil {
		res <- ErrInvalidURL
		return
	}

	content, _ := io.ReadAll(fh)
	defer fh.Close()

	if match, _ := regexp.Match(regexp.QuoteMeta(bookmark), content); match {
		res <- ErrDuplicate
		return
	}

	if titleFlag.GetValue() == titleFlag.GetDefault() {
		title, err := title(ctx, args[0])
		if err != nil {
			res <- err
			return
		}

		err = titleFlag.SetValue(title)
		if err != nil {
			res <- err
			return
		}
	}

	if titleFlag.GetValue() == titleFlag.GetDefault() {
		res <- ErrInvalidTitle
		return
	}

	_, err = fmt.Fprintf(fh, "%q %q\n", titleFlag.GetValue(), args[0])

	if err != nil {
		res <- err
	}

	close(res)
}

func title(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	page, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	res.Body.Close()

	return regex.MatchTitle(page), nil
}
