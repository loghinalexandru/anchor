package subcommand

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/loghinalexandru/anchor/internal/bookmark"
	"github.com/peterbourgon/ff/v4"
)

type get struct {
	command  ff.Command
	labels   []string
	fullFlag bool
	openFlag bool
}

func RegisterGet(root *ff.Command, rootFlags *ff.FlagSet) {
	cmd := get{}

	flags := ff.NewFlagSet("get").SetParent(rootFlags)
	_ = flags.StringSetVar(&cmd.labels, 'l', "label", "specify label hierarchy for each")
	_ = flags.BoolVar(&cmd.fullFlag, 'f', "full", "show full bookmark entry")
	_ = flags.BoolVar(&cmd.openFlag, 'o', "open", "open specified link")

	cmd.command = ff.Command{
		Name:      "get",
		Usage:     "get",
		ShortHelp: "get existing bookmarks",
		Flags:     flags,
		Exec: func(ctx context.Context, args []string) error {
			res := make(chan error, 1)
			go cmd.handle(args, res)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-res:
				return err
			}
		},
	}

	root.Subcommands = append(root.Subcommands, &cmd.command)
}

func (g *get) handle(args []string, res chan<- error) {
	defer close(res)

	rootDir, err := rootDir()
	if err != nil {
		res <- err
		return
	}

	err = validate(g.labels)
	if err != nil {
		res <- err
		return
	}

	path, err := os.Open(filepath.Join(rootDir, fileFrom(g.labels)))

	if err != nil {
		res <- err
		return
	}

	content, err := io.ReadAll(path)
	if err != nil {
		res <- err
		return
	}

	path.Close()

	var pattern string
	if len(args) >= 1 {
		pattern = args[0]
	}

	for _, l := range findLines(content, pattern) {
		title, url, err := bookmark.Parse(string(l))
		if err != nil {
			fmt.Print(url)
			res <- err
			return
		}

		if g.openFlag {
			err = open(url)
			if err != nil {
				res <- err
			}

			return
		}

		if g.fullFlag {
			fmt.Fprintln(os.Stdout, title, url)
		} else {
			fmt.Fprintln(os.Stdout, title)
		}
	}
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
