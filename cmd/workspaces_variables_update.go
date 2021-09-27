package cmd

import (
	"errors"
	"io"
)

type WorkspacesVariablesUpdateCmd struct {
	r    Runner
	deps dependencyProxies
	w    io.Writer
}

func newWorkspacesVariablesUpdateCmd(
	deps dependencyProxies,
	w io.Writer,
) *WorkspacesVariablesUpdateCmd {
	return &WorkspacesVariablesUpdateCmd{
		deps: deps,
		w:    w,
	}
}

func (c *WorkspacesVariablesUpdateCmd) Name() string {
	return "update"
}

func (c *WorkspacesVariablesUpdateCmd) Init(args []string) error {
	if len(args) < 1 {
		return errors.New("no subcommand given")
	}
	runners := []Runner{
		newWorkspacesVariablesUpdateValueCmd(c.deps, c.w),
		// newWorkspacesVariablesUpdateDescriptionCmd(c.deps, c.w),
	}
	return processSubcommand(&c.r, args, runners)
}

func (c *WorkspacesVariablesUpdateCmd) Run() error {
	return c.r.Run()
}
