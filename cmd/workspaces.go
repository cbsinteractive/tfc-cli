package cmd

import (
	"errors"
	"io"
)

type WorkspacesCmd struct {
	r       Runner
	deps    dependencyProxies
	w       io.Writer
	appName string
}

func NewWorkspacesCmd(options ExecuteOpts, deps dependencyProxies) *WorkspacesCmd {
	return &WorkspacesCmd{
		appName: options.AppName,
		deps:    deps,
		w:       options.Writer,
	}
}

func (c *WorkspacesCmd) Name() string {
	return "workspaces"
}

func (c *WorkspacesCmd) Init(args []string) error {
	if len(args) < 1 {
		return errors.New("no subcommand given")
	}
	runners := []Runner{
		newWorkspacesCreateCmd(c.deps, c.w, c.appName),
		newWorkspacesDeleteCmd(c.deps, c.w),
		newWorkspacesVariablesCmd(c.deps, c.w),
	}
	return processSubcommand(&c.r, args, runners)
}

func (c *WorkspacesCmd) Run() error {
	return c.r.Run()
}
