package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io"

	"github.com/hashicorp/go-tfe"
)

type VariableUpdateDescriptionOpts struct {
	key         string
	description string
}

type workspacesVariablesUpdateDescriptionCmd struct {
	fs   *flag.FlagSet
	deps dependencyProxies
	OrgOpts
	WorkspaceOpts
	VariableUpdateDescriptionOpts
	w io.Writer
}

type WorkspacesVariablesUpdateDescriptionCommandResult struct {
	Result *tfe.Variable
}

func newWorkspacesVariablesUpdateDescriptionCmd(deps dependencyProxies, w io.Writer) *workspacesVariablesUpdateDescriptionCmd {
	c := &workspacesVariablesUpdateDescriptionCmd{
		fs:   flag.NewFlagSet("description", flag.ContinueOnError),
		deps: deps,
		w:    w,
	}
	setCommonFlagsetOptions(c.fs, &c.OrgOpts, &c.WorkspaceOpts)
	c.fs.StringVar(&c.VariableUpdateDescriptionOpts.key, "key", "", string(VariableKeyUsage))
	c.fs.StringVar(&c.VariableUpdateDescriptionOpts.description, "description", "", string(VariableDescriptionUsage))
	return c
}

func (c *workspacesVariablesUpdateDescriptionCmd) Name() string {
	return c.fs.Name()
}

func (c *workspacesVariablesUpdateDescriptionCmd) Init(args []string) error {
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
	if c.VariableUpdateDescriptionOpts.key == "" {
		return errors.New("-key argument is required")
	}
	return nil
}

func (c *workspacesVariablesUpdateDescriptionCmd) Run() error {
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
	v, err := variableFromKey(client, c.deps.client, ctx, w.ID, c.VariableUpdateDescriptionOpts.key)
	if err != nil {
		return err
	}
	options := tfe.VariableUpdateOptions{
		Description: &c.VariableUpdateDescriptionOpts.description,
	}
	u, err := c.deps.client.workspacesCommands.variables.update(client, ctx, w.ID, v.ID, options)
	if err != nil {
		return err
	}
	if u == nil {
		return errors.New("variable and error both nil")
	}
	d, _ := json.Marshal(WorkspacesVariablesUpdateValueCommandResult{
		Result: u,
	})
	c.w.Write(d)
	return nil
}
