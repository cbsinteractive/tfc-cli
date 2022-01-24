package cmd

import (
	"context"
	"flag"
	"io"

	"github.com/hashicorp/go-tfe"
)

type workspacesUnsetDescriptionCmd struct {
	fs   *flag.FlagSet
	deps dependencyProxies
	OrgOpts
	WorkspaceOpts
	writer io.Writer
}

func newWorkspacesUnsetDescriptionCmd(deps dependencyProxies, w io.Writer) *workspacesUnsetDescriptionCmd {
	c := &workspacesUnsetDescriptionCmd{
		fs:     flag.NewFlagSet("unset-description", flag.ContinueOnError),
		deps:   deps,
		writer: w,
	}
	setCommonFlagsetOptions(c.fs, &c.OrgOpts, &c.WorkspaceOpts)
	return c
}

func (c *workspacesUnsetDescriptionCmd) Name() string {
	return c.fs.Name()
}

func (c *workspacesUnsetDescriptionCmd) Init(args []string) error {
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

func (c *workspacesUnsetDescriptionCmd) Run() error {
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
			Description: &emptyString,
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
