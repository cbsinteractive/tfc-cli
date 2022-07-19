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
	ID               string `json:"id"`
	Description      string `json:"description"`
	TerraformVersion string `json:"terraform-version"`
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
	c.appName = appName
	c.fs.StringVar(&c.WorkspaceOpts.terraformVersion, "terraformVersion", "", string(WorkspaceTerraformVersionUsage))
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
		Name:        &c.WorkspaceOpts.name,
		Description: &description,
	}
	if len(c.WorkspaceOpts.terraformVersion) > 0 {
		opts.TerraformVersion = &c.WorkspaceOpts.terraformVersion
	}
	workspace, err := c.deps.client.workspaces.create(
		client,
		context.Background(),
		c.OrgOpts.name,
		opts,
	)
	if err != nil {
		return err
	}
	if workspace == nil {
		return errors.New("workspace and error both nil")
	}
	output(c.w, WorkspacesCreateCommandResult{
		ID:               workspace.ID,
		Description:      workspace.Description,
		TerraformVersion: workspace.TerraformVersion,
	})
	return nil
}
