package subcommand

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/loghinalexandru/anchor/internal/bookmark"
	"github.com/peterbourgon/ff/v4"
)

type get struct {
	command *ff.Command
	labels  *[]string
}

func RegisterGet(root *ff.Command, rootFlags *ff.CoreFlags) {
	var g get
	var labels []string

	flags := ff.NewFlags("get").SetParent(rootFlags)
	_ = flags.StringSetVar(&labels, 'l', "label", "specify label hierarchy for each")
	_ = flags.Bool('f', "full", false, "show full bookmark entry")
	_ = flags.Bool('o', "open", false, "open specified link")

	g = get{
		command: &ff.Command{
			Name:      "get",
			Usage:     "get",
			ShortHelp: "get existing bookmarks",
			Flags:     flags,
			Exec: func(ctx context.Context, args []string) error {
				res := make(chan error, 1)
				go g.handle(args, res)

				select {
				case <-ctx.Done():
					return ctx.Err()
				case err := <-res:
					return err
				}
			},
		},
		labels: &labels,
	}

	root.Subcommands = append(root.Subcommands, g.command)
}

func (g get) handle(args []string, res chan<- error) {
	defer close(res)

	o, _ := g.command.Flags.GetFlag("open")
	f, _ := g.command.Flags.GetFlag("full")
	openFlag, _ := strconv.ParseBool(o.GetValue())
	fullFlag, _ := strconv.ParseBool(f.GetValue())
	dir, _ := g.command.Flags.GetFlag("root-dir")
	home, err := os.UserHomeDir()

	if err != nil {
		res <- err
		return
	}

	err = validate(*g.labels)
	if err != nil {
		res <- err
		return
	}

	tree := formatLabels(*g.labels)
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

		for _, l := range findLines(content, pattern) {
			title, url, err := bookmark.Parse(string(l))
			if err != nil {
				fmt.Print(url)
				res <- err
				return
			}

			if openFlag {
				err = open(url)
				if err != nil {
					res <- err
				}

				return
			}

			if fullFlag {
				fmt.Fprintln(os.Stdout, title, url)
			} else {
				fmt.Fprintln(os.Stdout, title)
			}
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
