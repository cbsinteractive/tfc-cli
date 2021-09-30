package cmd

import (
	"errors"
	"io"
)

type StateVersionsCurrentCmd struct {
	r    Runner
	deps dependencyProxies
	w    io.Writer
}

func NewStateVersionsCurrentCmd(deps dependencyProxies, w io.Writer) *StateVersionsCurrentCmd {
	return &StateVersionsCurrentCmd{
		deps: deps,
		w:    w,
	}
}

func (c *StateVersionsCurrentCmd) Name() string {
	return "current"
}

func (c *StateVersionsCurrentCmd) Init(args []string) error {
	if len(args) < 1 {
		err := errors.New("no subcommand given")
		c.w.Write(newCommandErrorOutput(err))
		return err
	}
	runners := []Runner{
		NewStateVersionsCurrentGetOutputCmd(c.deps, c.w),
	}
	return processSubcommand(&c.r, args, runners)
}

func (c *StateVersionsCurrentCmd) Run() error {
	return c.r.Run()
}
