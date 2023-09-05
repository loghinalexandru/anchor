package subcommand

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

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
	defer close(res)

	labelFlag, _ := c.Flags.GetFlag("label")
	dir, _ := c.Flags.GetFlag("root-dir")
	home, err := os.UserHomeDir()

	if err != nil {
		res <- err
		return
	}

	hierarchy := strings.Split(labelFlag.GetValue(), ",")
	path := fmt.Sprintf("%s/%s/%s", home, dir.GetValue(), strings.Join(hierarchy, "."))

	fh, err := os.Open(path)
	if err != nil {
		res <- err
		return
	}

	defer fh.Close()

	content, err := io.ReadAll(fh)
	if err != nil {
		res <- err
		return
	}

	if len(args) == 0 {
		fmt.Fprintf(os.Stdout, "%s\n", string(content))
		return
	}

	regex := regexp.MustCompile(fmt.Sprintf("(?im)^\".*%s.*\"$", regexp.QuoteMeta(args[0])))
	mm := regex.FindAll(content, -1)

	for _, m := range mm {
		fmt.Fprintf(os.Stdout, "%s\n", m)
	}
}
