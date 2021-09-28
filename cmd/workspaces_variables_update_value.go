package cmd

import (
	"context"
	"errors"
	"flag"
	"io"

	"github.com/hashicorp/go-tfe"
)

type WorkspacesVariablesUpdateValueCommandResult struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type VariableUpdateValueOpts struct {
	key   string
	value string
}

type workspacesVariablesUpdateValueCmd struct {
	fs   *flag.FlagSet
	deps dependencyProxies
	OrgOpts
	WorkspaceOpts
	VariableUpdateValueOpts
	w io.Writer
}

func newWorkspacesVariablesUpdateValueCmd(deps dependencyProxies, w io.Writer) *workspacesVariablesUpdateValueCmd {
	c := &workspacesVariablesUpdateValueCmd{
		fs:   flag.NewFlagSet("value", flag.ContinueOnError),
		deps: deps,
		w:    w,
	}
	setCommonFlagsetOptions(c.fs, &c.OrgOpts, &c.WorkspaceOpts)
	c.fs.StringVar(&c.VariableUpdateValueOpts.key, "key", "", string(VariableKeyUsage))
	c.fs.StringVar(&c.VariableUpdateValueOpts.value, "value", "", string(VariableValueUsage))
	return c
}

func (c *workspacesVariablesUpdateValueCmd) Name() string {
	return c.fs.Name()
}

func (c *workspacesVariablesUpdateValueCmd) Init(args []string) error {
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
	if c.VariableUpdateValueOpts.key == "" {
		return errors.New("-key argument is required")
	}
	return nil
}

func (c *workspacesVariablesUpdateValueCmd) Run() error {
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
	v, err := variableFromKey(client, c.deps.client, ctx, w.ID, c.VariableUpdateValueOpts.key)
	if err != nil {
		return err
	}
	options := tfe.VariableUpdateOptions{
		Value: &c.VariableUpdateValueOpts.value,
	}
	u, err := c.deps.client.workspacesCommands.variables.update(client, ctx, w.ID, v.ID, options)
	if err != nil {
		return err
	}
	if u == nil {
		return errors.New("variable and error both nil")
	}
	c.w.Write(newCommandResultOutput(WorkspacesVariablesUpdateValueCommandResult{
		ID:    u.ID,
		Key:   u.Key,
		Value: u.Value,
	}))
	return nil
}
