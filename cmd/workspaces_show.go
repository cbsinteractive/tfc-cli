package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io"

	"github.com/hashicorp/go-tfe"
)

type WorkspacesShowCommandResult struct {
	Result *tfe.Workspace `json:"result"`
}

type workspacesShowCmd struct {
	fs   *flag.FlagSet
	deps dependencyProxies
	OrgOpts
	WorkspaceOpts
	commonOpts
	w io.Writer
}

func newWorkspacesShowCmd(deps dependencyProxies, w io.Writer) *workspacesShowCmd {
	c := &workspacesShowCmd{
		fs:   flag.NewFlagSet("show", flag.ContinueOnError),
		deps: deps,
		w:    w,
	}
	setCommonFlagsetOptions(c.fs, &c.OrgOpts, &c.WorkspaceOpts)
	c.fs.BoolVar(&c.commonOpts.quiet, "quiet", false, string(QuietUsage))
	return c
}

func (c *workspacesShowCmd) Name() string {
	return c.fs.Name()
}

func (c *workspacesShowCmd) Init(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}
	if err := processCommonInputs(&c.OrgOpts.token, &c.OrgOpts.name, c.deps.os.lookupEnv); err != nil {
		return err
	}
	if c.WorkspaceOpts.name == "" {
		return errors.New("-workspace argument is required")
	}
	return nil
}

func (c *workspacesShowCmd) Run() error {
	client, err := tfe.NewClient(&tfe.Config{
		Token: c.OrgOpts.token,
	})
	if err != nil {
		return err
	}
	w, err := c.deps.client.workspaces.read(client, context.Background(), c.OrgOpts.name, c.WorkspaceOpts.name)
	if err != nil {
		return err
	}
	if w == nil {
		return errors.New("workspace and error both nil")
	}
	if c.commonOpts.quiet {
		return nil
	}
	d, _ := json.Marshal(WorkspacesShowCommandResult{
		Result: w,
	})
	c.w.Write(d)
	return nil
}
