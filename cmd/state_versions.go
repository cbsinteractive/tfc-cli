package cmd

import (
	"io"
)

type StateVersionsCmd struct {
	r    Runner
	deps dependencyProxies
	w    io.Writer
}

func NewStateVersionsCmd(deps dependencyProxies, w io.Writer) *StateVersionsCmd {
	return &StateVersionsCmd{
		deps: deps,
		w:    w,
	}
}

func (c *StateVersionsCmd) Name() string {
	return "stateversions"
}

func (c *StateVersionsCmd) Init(args []string) error {
	runners := []Runner{
		NewStateVersionsCurrentCmd(c.deps, c.w),
	}
	return processSubcommand(&c.r, args, runners)
}

func (c *StateVersionsCmd) Run() error {
	return c.r.Run()
}
