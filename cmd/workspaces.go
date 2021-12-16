package cmd

import (
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
		w:       options.Stdout,
	}
}

func (c *WorkspacesCmd) Name() string {
	return "workspaces"
}

func (c *WorkspacesCmd) Init(args []string) error {
	runners := []Runner{
		newWorkspacesCreateCmd(c.deps, c.w, c.appName),
		newWorkspacesDeleteCmd(c.deps, c.w),
		newWorkspacesShowCmd(c.deps, c.w),
		newWorkspacesUpdateCmd(c.deps, c.w),
		newWorkspacesVariablesCmd(c.deps, c.w),
		newWorkspacesVCSCmd(c.deps, c.w),
	}
	return processSubcommand(&c.r, args, runners)
}

func (c *WorkspacesCmd) Run() error {
	return c.r.Run()
}
