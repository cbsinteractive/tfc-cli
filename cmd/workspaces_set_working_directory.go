package cmd

import (
	"context"
	"flag"
	"io"

	"github.com/hashicorp/go-tfe"
)

type WorkspaceSetWorkingDirectoryOpts struct {
	workingDirectory string
}

type workspacesSetWorkingDirectoryCmd struct {
	fs   *flag.FlagSet
	deps dependencyProxies
	OrgOpts
	WorkspaceOpts
	WorkspaceSetWorkingDirectoryOpts
	writer io.Writer
}

func newWorkspacesSetWorkingDirectoryCmd(deps dependencyProxies, w io.Writer) *workspacesSetWorkingDirectoryCmd {
	c := &workspacesSetWorkingDirectoryCmd{
		fs:     flag.NewFlagSet("set-working-directory", flag.ContinueOnError),
		deps:   deps,
		writer: w,
	}
	setCommonFlagsetOptions(c.fs, &c.OrgOpts, &c.WorkspaceOpts)
	c.fs.StringVar(&c.WorkspaceSetWorkingDirectoryOpts.workingDirectory, "working-directory", "", string(WorkspaceWorkingDirectoryUsage))
	return c
}

func (c *workspacesSetWorkingDirectoryCmd) Name() string {
	return c.fs.Name()
}

func (c *workspacesSetWorkingDirectoryCmd) Init(args []string) error {
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

func (c *workspacesSetWorkingDirectoryCmd) Run() error {
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
			WorkingDirectory: &c.WorkspaceSetWorkingDirectoryOpts.workingDirectory,
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
