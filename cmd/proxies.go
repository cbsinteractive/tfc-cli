package cmd

import (
	"context"
	"os"

	"github.com/hashicorp/go-tfe"
)

type osProxy interface {
	lookupEnv(string) (string, bool)
}

type osProxyForProduction struct{}

func (p osProxyForProduction) lookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

func newOSProxy() osProxyForProduction {
	return osProxyForProduction{}
}

type stateVersionsProxy interface {
	currentWithOptions(*tfe.Client, context.Context, string, *tfe.StateVersionCurrentOptions) (*tfe.StateVersion, error)
}

type stateVersionsProxyForProduction struct{}

func newStateVersionsProxy() stateVersionsProxyForProduction {
	return stateVersionsProxyForProduction{}
}

func (c stateVersionsProxyForProduction) currentWithOptions(client *tfe.Client, ctx context.Context, workspaceID string, opts *tfe.StateVersionCurrentOptions) (*tfe.StateVersion, error) {
	return client.StateVersions.CurrentWithOptions(ctx, workspaceID, opts)
}

type workspacesProxy interface {
	create(*tfe.Client, context.Context, string, tfe.WorkspaceCreateOptions) (*tfe.Workspace, error)
	delete(*tfe.Client, context.Context, string, string) error
	read(*tfe.Client, context.Context, string, string) (*tfe.Workspace, error)
}

type workspacesProxyForProduction struct{}

func newWorkspacesProxy() workspacesProxyForProduction {
	return workspacesProxyForProduction{}
}

func (p workspacesProxyForProduction) create(client *tfe.Client, ctx context.Context, organization string, opts tfe.WorkspaceCreateOptions) (*tfe.Workspace, error) {
	return client.Workspaces.Create(ctx, organization, opts)
}

func (p workspacesProxyForProduction) delete(client *tfe.Client, ctx context.Context, organization string, workspace string) error {
	return client.Workspaces.Delete(ctx, organization, workspace)
}

func (p workspacesProxyForProduction) read(client *tfe.Client, ctx context.Context, organization string, workspace string) (*tfe.Workspace, error) {
	return client.Workspaces.Read(ctx, organization, workspace)
}

type clientProxy struct {
	stateVersions stateVersionsProxy
	workspaces    workspacesProxy
}

type dependencyProxies struct {
	os     osProxy
	client clientProxy
}
