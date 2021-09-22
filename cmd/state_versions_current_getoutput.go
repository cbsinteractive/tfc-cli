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

type OutputOpts struct {
	name string
}

type StateVersionsCurrentGetOutputCmd struct {
	fs   *flag.FlagSet
	deps dependencyCaller
	OrgOpts
	WorkspaceOpts
	OutputOpts
	w io.Writer
}

func NewStateVersionsCurrentGetOutputCmd(
	deps dependencyCaller,
	w io.Writer,
) *StateVersionsCurrentGetOutputCmd {
	c := &StateVersionsCurrentGetOutputCmd{
		fs:   flag.NewFlagSet("getoutput", flag.ContinueOnError),
		deps: deps,
		w:    w,
	}
	c.fs.StringVar(&c.OrgOpts.name, "org", "", string(OrgUsage))
	c.fs.StringVar(&c.OrgOpts.token, "token", "", string(TokenUsage))
	c.fs.StringVar(&c.WorkspaceOpts.name, "workspace", "", string(WorkspaceUsage))
	c.fs.StringVar(&c.OutputOpts.name, "name", "", string(OutputNameUsage))
	return c
}

func (c *StateVersionsCurrentGetOutputCmd) Name() string {
	return c.fs.Name()
}

func (c *StateVersionsCurrentGetOutputCmd) Init(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}
	if err := processCommonInputs(&c.OrgOpts.token, &c.OrgOpts.name, c.deps); err != nil {
		return err
	}
	if c.WorkspaceOpts.name == "" {
		return errors.New("-workspace argument is required")
	}
	if c.OutputOpts.name == "" {
		return errors.New("-name argument is required")
	}
	return nil
}

func (c *StateVersionsCurrentGetOutputCmd) Run() error {
	ctx := context.Background()
	client, err := tfe.NewClient(&tfe.Config{
		Token: c.OrgOpts.token,
	})
	if err != nil {
		return err
	}
	w, err := c.deps.clientWorkspacesRead(client, ctx, c.OrgOpts.name, c.WorkspaceOpts.name)
	if err != nil {
		return err
	}
	version, err := c.deps.clientStateVersionsCurrentWithOptions(
		client,
		ctx,
		w.ID,
		&tfe.StateVersionCurrentOptions{Include: "outputs"},
	)
	if err != nil {
		return err
	}
	for _, v := range version.Outputs {
		if v.Name == c.OutputOpts.name {
			d, _ := json.Marshal(CommandResult{
				Result: fmt.Sprintf("%v", v.Value),
			})
			c.w.Write(d)
			return nil
		}
	}
	d, _ := json.Marshal(CommandResult{
		Error: fmt.Sprintf("%s not found", c.OutputOpts.name),
	})
	c.w.Write(d)
	return nil
}
