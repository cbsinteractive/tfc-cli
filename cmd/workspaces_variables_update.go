package cmd

import (
	"errors"
	"flag"
	"io"
)

type WorkspacesVariablesUpdateCmd struct {
	fs   *flag.FlagSet
	deps dependencyProxies
	OrgOpts
	WorkspaceOpts
	w io.Writer
}

func newWorkspacesVariablesUpdateCmd(
	deps dependencyProxies,
	w io.Writer,
) *WorkspacesVariablesUpdateCmd {
	c := &WorkspacesVariablesUpdateCmd{
		fs:   flag.NewFlagSet("update", flag.ContinueOnError),
		deps: deps,
		w:    w,
	}
	setCommonFlagsetOptions(c.fs, c.OrgOpts, c.WorkspaceOpts)
	return c
}

func (c WorkspacesVariablesUpdateCmd) Name() string {
	return c.fs.Name()
}

func (c WorkspacesVariablesUpdateCmd) Init([]string) error {
	return errors.New("not implemented")
}

func (c WorkspacesVariablesUpdateCmd) Run() error {
	return errors.New("not implemented")
}
