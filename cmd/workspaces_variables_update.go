package cmd

import (
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
	fmt.Printf("Args: %+v\n", args)
	if err := c.fs.Parse(args); err != nil {
		return err
	}
	fmt.Printf("Command: %+v\n", c)
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

func (c *WorkspacesVariablesUpdateCmd) Run() error {
	// ctx := context.Background()
	_, err := tfe.NewClient(&tfe.Config{
		Token: c.OrgOpts.token,
	})
	if err != nil {
		return err
	}
	// Obtain the variable ID
	// v, err := c.deps.client.workspacesCommands.variables.read(client, ctx, c.VariableOpts.key)
	// if err != nil {
	// 	return err
	// }
	// c.VariableOpts.key
	return nil
}
