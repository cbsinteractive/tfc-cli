package cmd

import (
	"context"
	"flag"
	"io"

	"github.com/hashicorp/go-tfe"
)

type WorkspaceSetAutoApplyOpts struct {
	autoApply bool
}

type workspacesSetAutoApplyCmd struct {
	fs   *flag.FlagSet
	deps dependencyProxies
	OrgOpts
	WorkspaceOpts
	WorkspaceSetAutoApplyOpts
	writer io.Writer
}

func newWorkspacesSetAutoApplyCmd(deps dependencyProxies, w io.Writer) *workspacesSetAutoApplyCmd {
	c := &workspacesSetAutoApplyCmd{
		fs:     flag.NewFlagSet("set-auto-apply", flag.ContinueOnError),
		deps:   deps,
		writer: w,
	}
	setCommonFlagsetOptions(c.fs, &c.OrgOpts, &c.WorkspaceOpts)
	c.fs.BoolVar(&c.WorkspaceSetAutoApplyOpts.autoApply, "auto-apply", false, string(WorkspaceAutoApplyUsage))
	return c
}

func (c *workspacesSetAutoApplyCmd) Name() string {
	return c.fs.Name()
}

func (c *workspacesSetAutoApplyCmd) Init(args []string) error {
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
	return nil
}

func (c *workspacesSetAutoApplyCmd) Run() error {
	client, err := tfe.NewClient(&tfe.Config{
		Token: c.OrgOpts.token,
	})
	if err != nil {
		return err
	}
	updatedWorkspace, err := c.deps.client.workspaces.update(
		client,
		context.Background(),
		c.OrgOpts.name,
		c.WorkspaceOpts.name,
		tfe.WorkspaceUpdateOptions{
			AutoApply: &c.WorkspaceSetAutoApplyOpts.autoApply,
		},
	)
	if err != nil {
		return err
	}
	output(c.writer, WorkspacesUpdateCommandResult{
		ID:          updatedWorkspace.ID,
		Description: updatedWorkspace.Description,
	})
	return nil
}
