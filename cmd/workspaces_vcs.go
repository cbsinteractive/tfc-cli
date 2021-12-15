package cmd

import (
	"errors"
	"io"
)

type workspacesVCSCmd struct {
	r    Runner
	deps dependencyProxies
	w    io.Writer
}

func newWorkspacesVCSCmd(deps dependencyProxies, w io.Writer) *workspacesVCSCmd {
	return &workspacesVCSCmd{
		deps: deps,
		w:    w,
	}
}

func (c *workspacesVCSCmd) Name() string {
	return "vcs"
}

func (c *workspacesVCSCmd) Init(args []string) error {
	if len(args) < 1 {
		return errors.New("no subcommand given")
	}
	runners := []Runner{
		newWorkspacesVCSShowCmd(c.deps, c.w),
	}
	return processSubcommand(&c.r, args, runners)
}

func (c *workspacesVCSCmd) Run() error {
	return c.r.Run()
}
