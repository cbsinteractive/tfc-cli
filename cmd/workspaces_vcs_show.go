package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io"

	"github.com/hashicorp/go-tfe"
)

type workspacesVCSShowCmd struct {
	fs   *flag.FlagSet
	deps dependencyProxies
	OrgOpts
	WorkspaceOpts
	w io.Writer
}

func newWorkspacesVCSShowCmd(
	deps dependencyProxies,
	w io.Writer,
) *workspacesVCSShowCmd {
	c := &workspacesVCSShowCmd{
		fs:   flag.NewFlagSet("show", flag.ContinueOnError),
		deps: deps,
		w:    w,
	}
	setCommonFlagsetOptions(c.fs, &c.OrgOpts, &c.WorkspaceOpts)
	return c
}

func (c *workspacesVCSShowCmd) Name() string {
	return c.fs.Name()
}

func (c *workspacesVCSShowCmd) Init(args []string) error {
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

func (c *workspacesVCSShowCmd) Run() error {
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
	if w.VCSRepo == nil {
		d, _ := json.Marshal(CommandResult{
			Result: "VCS repo not set",
		})
		d = append(d, '\n')
		c.w.Write(d)
	}
	return nil
}
