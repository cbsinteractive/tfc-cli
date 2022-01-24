package cmd

import (
	"context"
	"flag"
	"io"

	"github.com/hashicorp/go-tfe"
)

type workspacesUnsetVCSBranchCmd struct {
	fs   *flag.FlagSet
	deps dependencyProxies
	OrgOpts
	WorkspaceOpts
	writer io.Writer
}

func newWorkspacesUnsetVCSBranchCmd(deps dependencyProxies, w io.Writer) *workspacesUnsetVCSBranchCmd {
	c := &workspacesUnsetVCSBranchCmd{
		fs:     flag.NewFlagSet("unset-vcs-branch", flag.ContinueOnError),
		deps:   deps,
		writer: w,
	}
	setCommonFlagsetOptions(c.fs, &c.OrgOpts, &c.WorkspaceOpts)
	return c
}

func (c *workspacesUnsetVCSBranchCmd) Name() string {
	return c.fs.Name()
}

func (c *workspacesUnsetVCSBranchCmd) Init(args []string) error {
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

func (c *workspacesUnsetVCSBranchCmd) Run() error {
	client, err := tfe.NewClient(&tfe.Config{
		Token: c.OrgOpts.token,
	})
	if err != nil {
		return err
	}
	updatedWorkspace, err := c.deps.client.workspaces.removeVCSConnection(
		client,
		context.Background(),
		c.OrgOpts.name,
		c.WorkspaceOpts.name,
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
