package cmd

import (
	"errors"
	"flag"
	"io"
)

type workspacesVCSShowCmd struct {
	fs   *flag.FlagSet
	deps dependencyProxies
	OrgOpts
	WorkspaceOpts
	w io.Writer
}

func newWorkspacesVCSShowCmd(
	deps dependencyProxies,
	w io.Writer,
) *workspacesVCSShowCmd {
	c := &workspacesVCSShowCmd{
		fs:   flag.NewFlagSet("show", flag.ContinueOnError),
		deps: deps,
		w:    w,
	}
	setCommonFlagsetOptions(c.fs, &c.OrgOpts, &c.WorkspaceOpts)
	return c
}

func (c *workspacesVCSShowCmd) Name() string {
	return c.fs.Name()
}

func (c *workspacesVCSShowCmd) Init(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}
	if err := processCommonInputs(
		&c.OrgOpts.token,
		&c.OrgOpts.name,
		c.deps.os.lookupEnv,
	); err != nil {
		return err
	}
	if c.WorkspaceOpts.name == "" {
		return errors.New("-workspace argument is required")
	}
	return nil
}

func (c *workspacesVCSShowCmd) Run() error {
	return nil
}
