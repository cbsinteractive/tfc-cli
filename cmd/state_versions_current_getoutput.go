package cmd

import (
	"context"
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
	deps dependencyProxies
	OrgOpts
	WorkspaceOpts
	OutputOpts
	w io.Writer
}

type StateVersionsCurrentGetOutputCommandResult struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

func NewStateVersionsCurrentGetOutputCmd(
	deps dependencyProxies,
	w io.Writer,
) *StateVersionsCurrentGetOutputCmd {
	c := &StateVersionsCurrentGetOutputCmd{
		fs:   flag.NewFlagSet("getoutput", flag.ContinueOnError),
		deps: deps,
		w:    w,
	}
	setCommonFlagsetOptions(c.fs, &c.OrgOpts, &c.WorkspaceOpts)
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
	w, err := c.deps.client.workspaces.read(client, ctx, c.OrgOpts.name, c.WorkspaceOpts.name)
	if err != nil {
		return err
	}
	version, err := c.deps.client.stateVersions.currentWithOptions(
		client,
		ctx,
		w.ID,
		&tfe.StateVersionCurrentOptions{Include: []tfe.StateVersionIncludeOpt{"outputs"}},
	)
	if err != nil {
		return err
	}
	for _, v := range version.Outputs {
		if v.Name == c.OutputOpts.name {
			output(c.w, StateVersionsCurrentGetOutputCommandResult{
				Name:  v.Name,
				Value: v.Value,
			})
			return nil
		}
	}
	return fmt.Errorf("%s not found", c.OutputOpts.name)
}
