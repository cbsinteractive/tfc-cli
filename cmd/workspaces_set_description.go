package cmd

import (
	"context"
	"flag"
	"io"

	"github.com/hashicorp/go-tfe"
)

type WorkspaceSetDescriptionOpts struct {
	description string
}

type workspacesSetDescriptionCmd struct {
	fs   *flag.FlagSet
	deps dependencyProxies
	OrgOpts
	WorkspaceOpts
	WorkspaceSetDescriptionOpts
	writer io.Writer
}

func newWorkspacesSetDescriptionCmd(deps dependencyProxies, w io.Writer) *workspacesSetDescriptionCmd {
	c := &workspacesSetDescriptionCmd{
		fs:     flag.NewFlagSet("set-description", flag.ContinueOnError),
		deps:   deps,
		writer: w,
	}
	setCommonFlagsetOptions(c.fs, &c.OrgOpts, &c.WorkspaceOpts)
	c.fs.StringVar(&c.WorkspaceSetDescriptionOpts.description, "description", "", string(WorkspaceDescriptionUsage))
	return c
}

func (c *workspacesSetDescriptionCmd) Name() string {
	return c.fs.Name()
}

func (c *workspacesSetDescriptionCmd) Init(args []string) error {
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

func (c *workspacesSetDescriptionCmd) Run() error {
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
			Description: &c.WorkspaceSetDescriptionOpts.description,
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
