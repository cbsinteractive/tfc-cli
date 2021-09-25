package cmd

import (
	"context"
	"errors"
	"flag"
	"io"

	"github.com/hashicorp/go-tfe"
)

type workspacesDeleteCmd struct {
	fs   *flag.FlagSet
	deps dependencyProxies
	OrgOpts
	WorkspaceOpts
	w io.Writer
}

func newWorkspacesDeleteCmd(deps dependencyProxies, w io.Writer) *workspacesDeleteCmd {
	c := &workspacesDeleteCmd{
		fs:   flag.NewFlagSet("delete", flag.ContinueOnError),
		deps: deps,
		w:    w,
	}
	c.fs.StringVar(&c.OrgOpts.name, "org", "", string(OrgUsage))
	c.fs.StringVar(&c.OrgOpts.token, "token", "", string(TokenUsage))
	c.fs.StringVar(&c.WorkspaceOpts.name, "workspace", "", string(WorkspaceUsage))
	return c
}

func (c *workspacesDeleteCmd) Name() string {
	return c.fs.Name()
}

func (c *workspacesDeleteCmd) Init(args []string) error {
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

func (c *workspacesDeleteCmd) Run() error {
	client, err := tfe.NewClient(&tfe.Config{
		Token: c.OrgOpts.token,
	})
	if err != nil {
		return err
	}
	err = c.deps.client.workspaces.delete(
		client,
		context.Background(),
		c.OrgOpts.name,
		c.WorkspaceOpts.name,
	)
	if err != nil {
		return err
	}
	return nil
}
