package command

import (
	"context"

	"github.com/peterbourgon/ff/v4"
)

const (
	initName      = "init"
	initUsage     = "anchor init [ARGS]"
	initShortHelp = "create a new empty home for anchor"
	initLongHelp  = `  Performs setup of the home directory depending on the selected
  backing storage via the config file.

  By default, local file system is used and there is no need
  for any other input. If you want to use something else provide
  what the backing storage requires as arguments.

EXAMPLES:
   # Setup local storage
   anchor init

   # Setup git storage
   anchor init git@github.com:loghinalexandru/anchor.git
`
)

type initCmd struct{}

func (init *initCmd) manifest(parent *ff.FlagSet) *ff.Command {
	return &ff.Command{
		Name:      initName,
		Usage:     initUsage,
		ShortHelp: initShortHelp,
		LongHelp:  initLongHelp,
		Flags:     ff.NewFlagSet("init").SetParent(parent),
		Exec: func(ctx context.Context, args []string) error {
			return init.handle(ctx.(rootContext), args)
		},
	}
}

func (init *initCmd) handle(ctx rootContext, args []string) error {
	return ctx.storer.Init(args...)
}
