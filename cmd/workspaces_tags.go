package cmd

import (
	"errors"
	"io"
)

type WorkspacesTagsCmd struct {
	r    Runner
	deps dependencyProxies
	w    io.Writer
}

func newWorkspacesTagsCmd(deps dependencyProxies, w io.Writer) *WorkspacesTagsCmd {
	return &WorkspacesTagsCmd{
		deps: deps,
		w:    w,
	}
}

func (c *WorkspacesTagsCmd) Name() string {
	return "tags"
}

func (c *WorkspacesTagsCmd) Init(args []string) error {
	if len(args) < 1 {
		return errors.New("no subcommand given")
	}
	runners := []Runner{
		newWorkspacesTagsCreateCmd(c.deps, c.w),
		// newWorkspacesTagsDeleteCmd(c.deps, c.w),
	}
	return processSubcommand(&c.r, args, runners)
}

func (c *WorkspacesTagsCmd) Run() error {
	return c.r.Run()
}
