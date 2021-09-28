package cmd

import (
	"context"
	"errors"
	"flag"
	"io"

	"github.com/hashicorp/go-tfe"
)

type VariableDeleteOpts struct {
	key string
}

type workspacesVariablesDeleteCmd struct {
	fs   *flag.FlagSet
	deps dependencyProxies
	OrgOpts
	WorkspaceOpts
	VariableDeleteOpts
	w io.Writer
}

func newWorkspacesVariablesDeleteCmd(deps dependencyProxies, w io.Writer) *workspacesVariablesDeleteCmd {
	c := &workspacesVariablesDeleteCmd{
		fs:   flag.NewFlagSet("delete", flag.ContinueOnError),
		deps: deps,
		w:    w,
	}
	setCommonFlagsetOptions(c.fs, &c.OrgOpts, &c.WorkspaceOpts)
	c.fs.StringVar(&c.VariableDeleteOpts.key, "key", "", string(VariableKeyUsage))
	return c
}

func (c *workspacesVariablesDeleteCmd) Name() string {
	return c.fs.Name()
}

func (c *workspacesVariablesDeleteCmd) Init(args []string) error {
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
	if c.WorkspaceOpts.name == "" {
		return errors.New("-workspace argument is required")
	}
	if c.VariableDeleteOpts.key == "" {
		return errors.New("-key argument is required")
	}
	return nil
}

func (c *workspacesVariablesDeleteCmd) Run() error {
	ctx := context.Background()
	client, err := tfe.NewClient(&tfe.Config{
		Token: c.OrgOpts.token,
	})
	if err != nil {
		return err
	}
	w, err := c.deps.client.workspaces.read(client, ctx, c.OrgOpts.name, c.WorkspaceOpts.name)
	if err != nil {
		return err
	}
	v, err := variableFromKey(client, c.deps.client, ctx, w.ID, c.VariableDeleteOpts.key)
	if err != nil {
		return err
	}
	err = c.deps.client.variables.delete(client, ctx, w.ID, v.ID)
	if err != nil {
		return err
	}
	return nil
}
