package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/hashicorp/go-tfe"
)

type WorkspacesCreateCommandResult struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type workspacesCreateCmd struct {
	fs   *flag.FlagSet
	deps dependencyProxies
	OrgOpts
	WorkspaceOpts
	w       io.Writer
	appName string
}

func newWorkspacesCreateCmd(deps dependencyProxies, w io.Writer, appName string) *workspacesCreateCmd {
	c := &workspacesCreateCmd{
		fs:   flag.NewFlagSet("create", flag.ContinueOnError),
		deps: deps,
		w:    w,
	}
	setCommonFlagsetOptions(c.fs, &c.OrgOpts, &c.WorkspaceOpts)
	c.fs.StringVar(&c.WorkspaceOpts.workingDirectory, "working-directory", "", string(WorkingDirectoryUsage))
	c.fs.StringVar(&c.WorkspaceOpts.vcsIdentifier, "vcs-identifier", "", string(VCSIdentifierUsage))
	c.fs.StringVar(&c.WorkspaceOpts.vcsBranch, "vcs-branch", "", string(VCSBranchUsage))
	c.fs.StringVar(&c.WorkspaceOpts.vcsOAuthTokenID, "vcs-oauth-token-id", "", string(VCSOAuthTokenIDUsage))
	c.appName = appName
	return c
}

func (c *workspacesCreateCmd) Name() string {
	return c.fs.Name()
}

func (c *workspacesCreateCmd) Init(args []string) error {
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

func (c *workspacesCreateCmd) Run() error {
	client, err := tfe.NewClient(&tfe.Config{
		Token: c.OrgOpts.token,
	})
	if err != nil {
		return err
	}
	description := fmt.Sprintf("Created by %s", c.appName)
	opts := tfe.WorkspaceCreateOptions{
		Name:             &c.WorkspaceOpts.name,
		Description:      &description,
		WorkingDirectory: &c.WorkspaceOpts.workingDirectory,
	}
	if c.WorkspaceOpts.vcsIdentifier != "" {
		if c.WorkspaceOpts.vcsBranch == "" {
			return errors.New("VCS identifier is specified but branch name is not")
		}
		if c.WorkspaceOpts.vcsOAuthTokenID == "" {
			return errors.New("VCS identifier is specified but OAuth token ID is not")
		}
		vcsOpts := tfe.VCSRepoOptions{
			Identifier:   &c.WorkspaceOpts.vcsIdentifier,
			Branch:       &c.WorkspaceOpts.vcsBranch,
			OAuthTokenID: &c.WorkspaceOpts.vcsOAuthTokenID,
		}
		opts.VCSRepo = &vcsOpts
	}
	w, err := c.deps.client.workspaces.create(
		client,
		context.Background(),
		c.OrgOpts.name,
		opts,
	)
	if err != nil {
		return err
	}
	if w == nil {
		return errors.New("workspace and error both nil")
	}
	output(c.w, WorkspacesCreateCommandResult{
		ID:          w.ID,
		Description: w.Description,
	})
	return nil
}
