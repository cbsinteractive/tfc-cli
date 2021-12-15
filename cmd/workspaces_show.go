package cmd

import (
	"context"
	"errors"
	"flag"
	"io"

	"github.com/hashicorp/go-tfe"
)

type WorkspacesShowCommandResult struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type workspacesShowCmd struct {
	fs   *flag.FlagSet
	deps dependencyProxies
	OrgOpts
	WorkspaceOpts
	commonOpts
	w io.Writer
}

func newWorkspacesShowCmd(deps dependencyProxies, w io.Writer) *workspacesShowCmd {
	c := &workspacesShowCmd{
		fs:   flag.NewFlagSet("show", flag.ContinueOnError),
		deps: deps,
		w:    w,
	}
	setCommonFlagsetOptions(c.fs, &c.OrgOpts, &c.WorkspaceOpts)
	c.fs.BoolVar(&c.commonOpts.quiet, "quiet", false, string(QuietUsage))
	return c
}

func (c *workspacesShowCmd) Name() string {
	return c.fs.Name()
}

func (c *workspacesShowCmd) Init(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}
	if err := processCommonInputs(&c.OrgOpts.token, &c.OrgOpts.name, c.deps.os.lookupEnv); err != nil {
		return err
	}
	if c.WorkspaceOpts.name == "" {
		return errors.New("-workspace argument is required")
	}
	return nil
}

func (c *workspacesShowCmd) Run() error {
	client, err := tfe.NewClient(&tfe.Config{
		Token: c.OrgOpts.token,
	})
	if err != nil {
		if !c.commonOpts.quiet {
			c.w.Write(newCommandErrorOutput(err))
		}
		return err
	}
	w, err := c.deps.client.workspaces.read(client, context.Background(), c.OrgOpts.name, c.WorkspaceOpts.name)
	if err != nil {
		return err
	}
	if w == nil {
		err := errors.New("workspace and error both nil")
		if !c.commonOpts.quiet {
			c.w.Write(newCommandErrorOutput(err))
		}
		return err
	}
	if c.commonOpts.quiet {
		return nil
	}
	c.w.Write(newCommandResultOutput(WorkspacesShowCommandResult{
		ID:          w.ID,
		Description: w.Description,
	}))
	return nil
}
