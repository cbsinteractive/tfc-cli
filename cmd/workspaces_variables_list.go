package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"strings"

	"github.com/hashicorp/go-tfe"
)

type WorkspacesVariablesListCmd struct {
	fs   *flag.FlagSet
	deps dependencyProxies
	OrgOpts
	WorkspaceOpts
	w io.Writer
}

func newWorkspacesVariablesListCmd(
	deps dependencyProxies,
	w io.Writer,
) *WorkspacesVariablesListCmd {
	c := &WorkspacesVariablesListCmd{
		fs:   flag.NewFlagSet("list", flag.ContinueOnError),
		deps: deps,
		w:    w,
	}
	setCommonFlagsetOptions(c.fs, &c.OrgOpts, &c.WorkspaceOpts)
	return c
}

func (c *WorkspacesVariablesListCmd) Name() string {
	return c.fs.Name()
}

func (c *WorkspacesVariablesListCmd) Init(args []string) error {
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
	return nil
}

type WorkspacesVariablesListCommandResult struct {
	Result string
}

func (c *WorkspacesVariablesListCmd) Run() error {
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
	l, err := c.deps.client.workspacesCommands.variables.list(client, ctx, w.ID, &tfe.VariableListOptions{})
	if err != nil {
		return err
	}
	keys := make([]string, 0)
	for _, i := range l.Items {
		keys = append(keys, i.Key)
	}
	d, _ := json.Marshal(WorkspacesVariablesListCommandResult{
		Result: strings.Join(keys, ","),
	})
	c.w.Write(d)
	return nil
}
