package cmd

import (
	"context"
	"flag"
	"io"

	"github.com/hashicorp/go-tfe"
)

type workspacesUnsetWorkingDirectoryCmd struct {
	fs   *flag.FlagSet
	deps dependencyProxies
	OrgOpts
	WorkspaceOpts
	writer io.Writer
}

func newWorkspacesUnsetWorkingDirectoryCmd(deps dependencyProxies, w io.Writer) *workspacesUnsetWorkingDirectoryCmd {
	c := &workspacesUnsetWorkingDirectoryCmd{
		fs:     flag.NewFlagSet("unset-working-directory", flag.ContinueOnError),
		deps:   deps,
		writer: w,
	}
	setCommonFlagsetOptions(c.fs, &c.OrgOpts, &c.WorkspaceOpts)
	return c
}

func (c *workspacesUnsetWorkingDirectoryCmd) Name() string {
	return c.fs.Name()
}

func (c *workspacesUnsetWorkingDirectoryCmd) Init(args []string) error {
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

func (c *workspacesUnsetWorkingDirectoryCmd) Run() error {
	client, err := tfe.NewClient(&tfe.Config{
		Token: c.OrgOpts.token,
	})
	if err != nil {
		return err
	}
	emptyString := ""
	updatedWorkspace, err := c.deps.client.workspaces.update(
		client,
		context.Background(),
		c.OrgOpts.name,
		c.WorkspaceOpts.name,
		tfe.WorkspaceUpdateOptions{
			WorkingDirectory: &emptyString,
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
