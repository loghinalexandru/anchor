package subcommand

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/loghinalexandru/anchor/internal/regex"
	"github.com/peterbourgon/ff/v4"
)

type getCmd ff.Command

func RegisterGet(root *ff.Command, rootFlags *ff.CoreFlags) {
	var cmd *getCmd
	flags := ff.NewFlags("get").SetParent(rootFlags)
	_ = flags.String('l', "label", "root", "specify label hierarchy for each")

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
	labelFlag, _ := c.Flags.GetFlag("label")
	dir, _ := c.Flags.GetFlag("root-dir")
	home, err := os.UserHomeDir()

	if err != nil {
		res <- err
		return
	}

	paths := multiLevelPaths(filepath.Join(home, dir.GetValue()), labelFlag)

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

		if len(args) == 0 {
			fmt.Fprintf(os.Stdout, "%s\n", string(content))
			continue
		}

		for _, m := range regex.MatchLines(content, args[0]) {
			fmt.Fprintf(os.Stdout, "%s\n", m)
		}
	}

	close(res)
}

func multiLevelPaths(rootDir string, labels ff.Flag) []string {
	var paths []string

	if labels.GetValue() == labels.GetDefault() {
		return []string{filepath.Join(rootDir, labels.GetDefault())}
	}

	ll := strings.Split(labels.GetValue(), ",")

	for i := 1; i <= len(ll); i++ {
		paths = append(paths, filepath.Join(rootDir, strings.Join(ll[:i], ".")))
	}

	return paths
}
