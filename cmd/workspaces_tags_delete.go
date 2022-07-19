package cmd

import (
	"context"
	"errors"
	"flag"
	"io"

	"github.com/hashicorp/go-tfe"
)

type TagDeleteOpts struct {
	tag string
}

type workspacesTagsDeleteCmd struct {
	fs   *flag.FlagSet
	deps dependencyProxies
	OrgOpts
	WorkspaceOpts
	TagDeleteOpts
	w io.Writer
}

func newWorkspacesTagsDeleteCmd(deps dependencyProxies, w io.Writer) *workspacesTagsDeleteCmd {
	c := &workspacesTagsDeleteCmd{
		fs:   flag.NewFlagSet("delete", flag.ContinueOnError),
		deps: deps,
		w:    w,
	}
	setCommonFlagsetOptions(c.fs, &c.OrgOpts, &c.WorkspaceOpts)
	c.fs.StringVar(&c.TagDeleteOpts.tag, "tag", "", string(WorkspaceTagUsage))
	return c
}

func (c *workspacesTagsDeleteCmd) Name() string {
	return c.fs.Name()
}

func (c *workspacesTagsDeleteCmd) Init(args []string) error {
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
	if c.TagDeleteOpts.tag == "" {
		return errors.New("-tag argument is required")
	}
	return nil
}

func (c *workspacesTagsDeleteCmd) Run() error {
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
	options := tfe.WorkspaceRemoveTagsOptions{Tags: []*tfe.Tag{
		{
			Name: c.TagDeleteOpts.tag,
		},
	}}
	err = c.deps.client.tags.delete(client, ctx, w.ID, options)
	if err != nil {
		return err
	}
	return nil
}
