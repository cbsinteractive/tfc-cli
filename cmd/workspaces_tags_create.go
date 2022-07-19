package cmd

import (
	"context"
	"errors"
	"flag"
	"io"

	"github.com/hashicorp/go-tfe"
)

type TagCreateOpts struct {
	tag string
}

type workspacesTagsCreateCmd struct {
	fs   *flag.FlagSet
	deps dependencyProxies
	OrgOpts
	WorkspaceOpts
	TagCreateOpts
	w io.Writer
}

func newWorkspacesTagsCreateCmd(deps dependencyProxies, w io.Writer) *workspacesTagsCreateCmd {
	c := &workspacesTagsCreateCmd{
		fs:   flag.NewFlagSet("create", flag.ContinueOnError),
		deps: deps,
		w:    w,
	}
	setCommonFlagsetOptions(c.fs, &c.OrgOpts, &c.WorkspaceOpts)
	c.fs.StringVar(&c.TagCreateOpts.tag, "tag", "", string(WorkspaceTagUsage))
	return c
}

func (c *workspacesTagsCreateCmd) Name() string {
	return c.fs.Name()
}

func (c *workspacesTagsCreateCmd) Init(args []string) error {
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
	if c.TagCreateOpts.tag == "" {
		return errors.New("-tag argument is required")
	}
	return nil
}

func (c *workspacesTagsCreateCmd) Run() error {
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
	options := tfe.WorkspaceAddTagsOptions{Tags: []*tfe.Tag{
		{
			Name: c.TagCreateOpts.tag,
		},
	}}
	err = c.deps.client.tags.create(client, ctx, w.ID, options)
	if err != nil {
		return err
	}
	return nil
}
