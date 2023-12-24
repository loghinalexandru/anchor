package command

import (
	"context"
	"fmt"
	"runtime"

	"github.com/peterbourgon/ff/v4"
)

const (
	versionName  = "version"
	versionUsage = "anchor version"
	versionShort = "print anchor version"
)

var (
	// version contains latest semver value set via -ldflags
	// and should not be changed by hand.
	version = "dev"
)

type versionCmd struct{}

func (init *versionCmd) manifest(parent *ff.FlagSet) *ff.Command {
	return &ff.Command{
		Name:      versionName,
		Usage:     versionUsage,
		ShortHelp: versionShort,
		Flags:     ff.NewFlagSet("version").SetParent(parent),
		Exec: func(ctx context.Context, args []string) error {
			fmt.Printf("anchor version %s %s/%s\n", version, runtime.GOOS, runtime.GOARCH)
			return nil
		},
	}
}
