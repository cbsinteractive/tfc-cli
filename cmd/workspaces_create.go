package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/hashicorp/go-tfe"
)

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
	return c
}

func (c *workspacesCreateCmd) Name() string {
	return c.fs.Name()
}

func (c *workspacesCreateCmd) Init(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		c.w.Write(newCommandErrorOutput(err))
		return err
	}
	if err := processCommonInputs(&c.OrgOpts.token, &c.OrgOpts.name, c.deps.os.lookupEnv); err != nil {
		c.w.Write(newCommandErrorOutput(err))
		return err
	}
	if c.WorkspaceOpts.name == "" {
		err := errors.New("-workspace argument is required")
		c.w.Write(newCommandErrorOutput(err))
		return err
	}
	return nil
}

func (c *workspacesCreateCmd) Run() error {
	client, err := tfe.NewClient(&tfe.Config{
		Token: c.OrgOpts.token,
	})
	if err != nil {
		c.w.Write(newCommandErrorOutput(err))
		return err
	}
	description := fmt.Sprintf("Created by %s", c.appName)
	workspace, err := c.deps.client.workspaces.create(client, context.Background(), c.OrgOpts.name, tfe.WorkspaceCreateOptions{
		Name:        &c.WorkspaceOpts.name,
		Description: &description,
	})
	if err != nil {
		c.w.Write(newCommandErrorOutput(err))
		return err
	}
	if workspace == nil {
		err := errors.New("workspace and error both nil")
		c.w.Write(newCommandErrorOutput(err))
		return err
	}
	return nil
}
