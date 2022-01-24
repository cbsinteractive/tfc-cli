package cmd

import (
	"context"
	"errors"
	"flag"
	"io"

	"github.com/hashicorp/go-tfe"
)

type WorkspaceSetVCSBranchOpts struct {
	identifier   string
	branch       string
	oAuthTokenID string
}

type workspacesSetVCSBranchCmd struct {
	fs   *flag.FlagSet
	deps dependencyProxies
	OrgOpts
	WorkspaceOpts
	WorkspaceSetVCSBranchOpts
	writer io.Writer
}

func newWorkspacesSetVCSBranchCmd(deps dependencyProxies, w io.Writer) *workspacesSetVCSBranchCmd {
	c := &workspacesSetVCSBranchCmd{
		fs:     flag.NewFlagSet("set-vcs-branch", flag.ContinueOnError),
		deps:   deps,
		writer: w,
	}
	setCommonFlagsetOptions(c.fs, &c.OrgOpts, &c.WorkspaceOpts)
	c.fs.StringVar(&c.WorkspaceSetVCSBranchOpts.branch, "branch", "", string(VCSBranchUsage))
	c.fs.StringVar(&c.WorkspaceSetVCSBranchOpts.identifier, "identifier", "", string(VCSIdentifierUsage))
	c.fs.StringVar(&c.WorkspaceSetVCSBranchOpts.oAuthTokenID, "oauth-token-id", "", string(VCSOAuthTokenIDUsage))
	return c
}

func (c *workspacesSetVCSBranchCmd) Name() string {
	return c.fs.Name()
}

func (c *workspacesSetVCSBranchCmd) Init(args []string) error {
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
	if c.WorkspaceSetVCSBranchOpts.identifier == "" {
		return errors.New("-identifier argument is required")
	}
	if c.WorkspaceSetVCSBranchOpts.branch == "" {
		return errors.New("-branch argument is required")
	}
	if c.WorkspaceSetVCSBranchOpts.oAuthTokenID == "" {
		return errors.New("-oauth-token-id argument is required")
	}
	return nil
}

func (c *workspacesSetVCSBranchCmd) Run() error {
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
			VCSRepo: &tfe.VCSRepoOptions{
				Identifier:   &c.WorkspaceSetVCSBranchOpts.identifier,
				Branch:       &c.WorkspaceSetVCSBranchOpts.branch,
				OAuthTokenID: &c.WorkspaceSetVCSBranchOpts.oAuthTokenID,
			},
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
