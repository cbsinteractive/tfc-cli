package cmd

import (
	"context"
	"errors"
	"flag"
	"io"

	"github.com/hashicorp/go-tfe"
)

type WorkspacesUpdateCommandResult struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type workspacesUpdateCmd struct {
	flagSet *flag.FlagSet
	deps    dependencyProxies
	OrgOpts
	WorkspaceOpts
	writer io.Writer
}

func newWorkspacesUpdateCmd(deps dependencyProxies, writer io.Writer) *workspacesUpdateCmd {
	c := &workspacesUpdateCmd{
		flagSet: flag.NewFlagSet("update", flag.ContinueOnError),
		deps:    deps,
		writer:  writer,
	}
	setCommonFlagsetOptions(c.flagSet, &c.OrgOpts, &c.WorkspaceOpts)
	c.flagSet.StringVar(&c.WorkspaceOpts.description, "description", "", string(WorkspaceDescriptionUsage))
	c.flagSet.StringVar(&c.WorkspaceOpts.workingDirectory, "working-directory", "", string(WorkingDirectoryUsage))
	c.flagSet.StringVar(&c.WorkspaceOpts.vcsIdentifier, "vcs-identifier", "", string(VCSIdentifierUsage))
	c.flagSet.StringVar(&c.WorkspaceOpts.vcsBranch, "vcs-branch", "", string(VCSBranchUsage))
	c.flagSet.StringVar(&c.WorkspaceOpts.vcsOAuthTokenID, "vcs-oauth-token-id", "", string(VCSOAuthTokenIDUsage))
	return c
}

func (c *workspacesUpdateCmd) Name() string {
	return c.flagSet.Name()
}

func (c *workspacesUpdateCmd) Init(args []string) error {
	if err := c.flagSet.Parse(args); err != nil {
		return err
	}
	if err := processCommonInputs(&c.OrgOpts.token, &c.OrgOpts.name, c.deps.os.lookupEnv); err != nil {
		return err
	}
	if c.WorkspaceOpts.name == "" {
		return errors.New("-workspace argument is required")
	}
	if c.WorkspaceOpts.vcsIdentifier != "" {
		if c.WorkspaceOpts.vcsBranch == "" {
			return errors.New("VCS identifier is specified but branch name is not")
		}
		if c.WorkspaceOpts.vcsOAuthTokenID == "" {
			return errors.New("VCS identifier is specified but OAuth token ID is not")
		}
	}
	return nil
}

func (c *workspacesUpdateCmd) Run() error {
	client, err := tfe.NewClient(&tfe.Config{
		Token: c.OrgOpts.token,
	})
	if err != nil {
		return err
	}
	existingWorkspace, err := c.deps.client.workspaces.read(
		client,
		context.Background(),
		c.OrgOpts.name,
		c.WorkspaceOpts.name,
	)
	if err != nil {
		return err
	}

	// Initialize update options from those we just read
	updateOpts := tfe.WorkspaceUpdateOptions{
		Description: &existingWorkspace.Description,
	}
	if c.WorkspaceOpts.description != "" {
		updateOpts.Description = &c.WorkspaceOpts.description
	}

	updatedWorkspace, err := c.deps.client.workspaces.update(
		client,
		context.Background(),
		c.OrgOpts.name,
		c.WorkspaceOpts.name,
		updateOpts,
	)
	if err != nil {
		return nil
	}
	output(c.writer, WorkspacesUpdateCommandResult{
		ID:          updatedWorkspace.ID,
		Description: updatedWorkspace.Description,
	})
	return nil
}
