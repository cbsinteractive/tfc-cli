package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/hashicorp/go-tfe"
)

type WorkspacesVariablesUpdateCmd struct {
	fs   *flag.FlagSet
	deps dependencyProxies
	OrgOpts
	WorkspaceOpts
	VariableOpts
	w io.Writer
}

func newWorkspacesVariablesUpdateCmd(
	deps dependencyProxies,
	w io.Writer,
) *WorkspacesVariablesUpdateCmd {
	c := &WorkspacesVariablesUpdateCmd{
		fs:   flag.NewFlagSet("update", flag.ContinueOnError),
		deps: deps,
		w:    w,
	}
	setCommonFlagsetOptions(c.fs, &c.OrgOpts, &c.WorkspaceOpts)
	c.fs.StringVar(&c.VariableOpts.key, "key", "", string(VariableKeyUsage))
	c.fs.StringVar(&c.VariableOpts.value, "value", "", string(VariableValueUsage))
	c.fs.StringVar(&c.VariableOpts.description, "description", "", string(VariableDescriptionUsage))
	c.fs.StringVar(&c.VariableOpts.category, "category", "", string(VariableCategoryUsage))
	c.fs.BoolVar(&c.VariableOpts.hcl, "hcl", false, string(VariableHCLUsage))
	c.fs.BoolVar(&c.VariableOpts.sensitive, "sensitive", false, string(VariableSensitiveUsage))
	return c
}

func (c *WorkspacesVariablesUpdateCmd) Name() string {
	return c.fs.Name()
}

func (c *WorkspacesVariablesUpdateCmd) Init(args []string) error {
	fmt.Printf("Args: %v\n", args)
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
	if c.VariableOpts.key == "" {
		return errors.New("-key argument is required")
	}
	return nil
}

func variableFromKey(client *tfe.Client, proxy clientProxy, ctx context.Context, workspaceID string, key string) (*tfe.Variable, error) {
	v, err := proxy.workspacesCommands.variables.list(client, ctx, workspaceID, tfe.VariableListOptions{})
	if err != nil {
		return nil, err
	}
	for _, i := range v.Items {
		if i.Key == key {
			return i, nil
		}
	}
	return nil, fmt.Errorf("variable %s not found", key)
}

type WorkspacesVariablesUpdateCommandResult struct {
	Result *tfe.Variable
}

func (c *WorkspacesVariablesUpdateCmd) Run() error {
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
	v, err := variableFromKey(client, c.deps.client, ctx, w.ID, c.VariableOpts.key)
	if err != nil {
		return err
	}
	fmt.Printf("Variable opts: %+v\n", c.VariableOpts)
	options := tfe.VariableUpdateOptions{
		Value:       &c.VariableOpts.value,
		Description: &c.VariableOpts.description,
	}
	fmt.Printf("Description: %s\n", *options.Description)
	u, err := c.deps.client.workspacesCommands.variables.update(client, ctx, w.ID, v.ID, options)
	if err != nil {
		return err
	}
	if u == nil {
		return errors.New("variable and error both nil")
	}
	d, _ := json.Marshal(WorkspacesVariablesUpdateCommandResult{
		Result: u,
	})
	c.w.Write(d)
	return nil
}
