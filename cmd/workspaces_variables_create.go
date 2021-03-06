package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/hashicorp/go-tfe"
)

type VariableCreateOpts struct {
	key         string
	value       string
	description string
	categoryRaw string
	category    tfe.CategoryType
	sensitive   bool
	hcl         bool
}

type workspacesVariablesCreateCmd struct {
	fs   *flag.FlagSet
	deps dependencyProxies
	OrgOpts
	WorkspaceOpts
	VariableCreateOpts
	w io.Writer
}

type WorkspacesVariablesCreateCommandResult struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Value       string `json:"value"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Sensitive   bool   `json:"sensitive"`
	HCL         bool   `json:"hcl"`
}

func newWorkspacesVariablesCreateCmd(deps dependencyProxies, w io.Writer) *workspacesVariablesCreateCmd {
	c := &workspacesVariablesCreateCmd{
		fs:   flag.NewFlagSet("create", flag.ContinueOnError),
		deps: deps,
		w:    w,
	}
	setCommonFlagsetOptions(c.fs, &c.OrgOpts, &c.WorkspaceOpts)
	c.fs.StringVar(&c.VariableCreateOpts.key, "key", "", string(VariableKeyUsage))
	c.fs.StringVar(&c.VariableCreateOpts.value, "value", "", string(VariableValueUsage))
	c.fs.StringVar(&c.VariableCreateOpts.categoryRaw, "category", "", string(VariableCategoryUsage))
	c.fs.StringVar(&c.VariableCreateOpts.description, "description", "", string(VariableDescriptionUsage))
	c.fs.BoolVar(&c.VariableCreateOpts.sensitive, "sensitive", false, string(VariableSensitiveUsage))
	c.fs.BoolVar(&c.VariableCreateOpts.hcl, "hcl", false, string(VariableHCLUsage))
	return c
}

func (c *workspacesVariablesCreateCmd) Name() string {
	return c.fs.Name()
}

func (c *workspacesVariablesCreateCmd) Init(args []string) error {
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
	if c.VariableCreateOpts.key == "" {
		return errors.New("-key argument is required")
	}
	switch c.VariableCreateOpts.categoryRaw {
	case "terraform":
		c.VariableCreateOpts.category = tfe.CategoryTerraform
	case "env":
		c.VariableCreateOpts.category = tfe.CategoryEnv
	default:
		return fmt.Errorf(`invalid category: "%s"`, c.VariableCreateOpts.categoryRaw)
	}
	return nil
}

func (c *workspacesVariablesCreateCmd) Run() error {
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
	options := tfe.VariableCreateOptions{
		Key:         &c.VariableCreateOpts.key,
		Value:       &c.VariableCreateOpts.value,
		Category:    &c.VariableCreateOpts.category,
		Description: &c.VariableCreateOpts.description,
		Sensitive:   &c.VariableCreateOpts.sensitive,
		HCL:         &c.VariableCreateOpts.hcl,
	}
	v, err := c.deps.client.variables.create(client, ctx, w.ID, options)
	if err != nil {
		return err
	}
	if v == nil {
		return errors.New("variable and error both nil")
	}
	output(c.w, WorkspacesVariablesCreateCommandResult{
		ID:  v.ID,
		Key: v.Key,
		// The API provides an empty string if the variable is set as sensitive
		Value:       v.Value,
		Category:    string(v.Category),
		Description: v.Description,
		Sensitive:   v.Sensitive,
		HCL:         v.HCL,
	})
	return nil
}
