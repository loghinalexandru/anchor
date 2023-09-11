package subcommand

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/loghinalexandru/anchor/internal/bookmark"
	"github.com/loghinalexandru/anchor/internal/regex"
	"github.com/peterbourgon/ff/v4"
)

type getCmd ff.Command

func RegisterGet(root *ff.Command, rootFlags *ff.CoreFlags) {
	var cmd *getCmd
	flags := ff.NewFlags("get").SetParent(rootFlags)
	_ = flags.StringSet('l', "label", "specify label hierarchy for each")
	_ = flags.Bool('o', "open", false, "open specified link")

	cmd = &getCmd{
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

	root.Subcommands = append(root.Subcommands, (*ff.Command)(cmd))
}

func (c *getCmd) handle(args []string, res chan<- error) {
	defer close(res)

	labelFlag, _ := c.Flags.GetFlag("label")
	openFlag, _ := c.Flags.GetFlag("open")
	dir, _ := c.Flags.GetFlag("root-dir")
	home, err := os.UserHomeDir()

	if err != nil {
		res <- err
		return
	}

	tree, err := flattenWithValidation(labelFlag.GetValue())
	if err != nil {
		res <- err
		return
	}

	paths, err := multiLevelPaths(filepath.Join(home, dir.GetValue()), tree)
	if err != nil {
		res <- err
		return
	}

	for _, p := range paths {
		fh, err := os.Open(p)
		if err != nil {
			res <- err
			return
		}

		content, err := io.ReadAll(fh)
		if err != nil {
			res <- err
			return
		}

		fh.Close()

		var pattern string
		if len(args) >= 1 {
			pattern = args[0]
		}

		for _, l := range regex.FindLines(content, pattern) {
			title, url, err := bookmark.Parse(string(l))
			if err != nil {
				fmt.Print(url)
				res <- err
				return
			}

			if openFlag.GetValue() == "true" {
				err = open(url)
				if err != nil {
					res <- err
				}

				return
			}

			fmt.Fprintln(os.Stdout, title)
		}
	}
}

func multiLevelPaths(rootDir string, treePath string) ([]string, error) {
	var paths []string

	dd, err := os.ReadDir(rootDir)
	if err != nil {
		return nil, err
	}

	for _, d := range dd {
		if d.IsDir() {
			continue
		}

		if strings.HasPrefix(d.Name(), treePath) {
			paths = append(paths, filepath.Join(rootDir, d.Name()))
		}
	}

	return paths, nil
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
