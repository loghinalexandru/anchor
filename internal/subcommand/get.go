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

type getCmd struct {
	command  ff.Command
	labels   []string
	fullFlag bool
	openFlag bool
}

func RegisterGet(root *ff.Command, rootFlags *ff.FlagSet) {
	cmd := getCmd{}

	flags := ff.NewFlagSet("get").SetParent(rootFlags)
	_ = flags.StringSetVar(&cmd.labels, 'l', "label", "specify label hierarchy for each")
	_ = flags.BoolVar(&cmd.fullFlag, 'f', "full", "show full bookmark entry")
	_ = flags.BoolVar(&cmd.openFlag, 'o', "open", "open specified link")

	cmd.command = ff.Command{
		Name:      "get",
		Usage:     "get",
		ShortHelp: "get existing bookmarks",
		Flags:     flags,
		Exec:      handlerMiddleware(cmd.handle),
	}

	root.Subcommands = append(root.Subcommands, &cmd.command)
}

func (get *getCmd) handle(_ context.Context, args []string) error {

	dir, err := rootDir()
	if err != nil {
		return err
	}

	err = validate(get.labels)
	if err != nil {
		return err
	}

	path, err := os.Open(filepath.Join(dir, fileFrom(get.labels)))

	if err != nil {
		return err
	}

	content, err := io.ReadAll(path)
	if err != nil {
		return err
	}

	_ = path.Close()
	var pattern string
	if len(args) >= 1 {
		pattern = args[0]
	}

	for _, l := range findLines(content, pattern) {
		title, url, err := bookmark.Parse(string(l))
		if err != nil {
			fmt.Print(url)
			return err
		}

		if get.openFlag {
			err = open(url)
			if err != nil {
				return err
			}

			return nil
		}

		if get.fullFlag {
			_, _ = fmt.Fprintln(os.Stdout, title, url)
		} else {
			_, _ = fmt.Fprintln(os.Stdout, title)
		}
	}

	return nil
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
