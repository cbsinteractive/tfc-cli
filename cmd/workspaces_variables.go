package cmd

import (
	"errors"
	"io"
)

type WorkspacesVariablesCmd struct {
	r    Runner
	deps dependencyProxies
	w    io.Writer
}

func newWorkspacesVariablesCmd(deps dependencyProxies, w io.Writer) *WorkspacesVariablesCmd {
	return &WorkspacesVariablesCmd{
		deps: deps,
		w:    w,
	}
}

func (c *WorkspacesVariablesCmd) Name() string {
	return "variables"
}

func (c *WorkspacesVariablesCmd) Init(args []string) error {
	if len(args) < 1 {
		return errors.New("no subcommand given")
	}
	runners := []Runner{
		newWorkspacesVariablesListCmd(c.deps, c.w),
		newWorkspacesVariablesUpdateCmd(c.deps, c.w),
	}
	return processSubcommand(&c.r, args, runners)
}

func (c *WorkspacesVariablesCmd) Run() error {
	return c.r.Run()
}
