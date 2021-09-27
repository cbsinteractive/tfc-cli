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
	currentWithOptions(client *tfe.Client, ctx context.Context, workspaceID string, options *tfe.StateVersionCurrentOptions) (*tfe.StateVersion, error)
}

type stateVersionsProxyForProduction struct{}

func newStateVersionsProxy() stateVersionsProxyForProduction {
	return stateVersionsProxyForProduction{}
}

func (c stateVersionsProxyForProduction) currentWithOptions(client *tfe.Client, ctx context.Context, workspaceID string, options *tfe.StateVersionCurrentOptions) (*tfe.StateVersion, error) {
	return client.StateVersions.CurrentWithOptions(ctx, workspaceID, options)
}

type workspacesProxy interface {
	create(client *tfe.Client, ctx context.Context, organization string, options tfe.WorkspaceCreateOptions) (*tfe.Workspace, error)
	delete(client *tfe.Client, ctx context.Context, organization string, workspace string) error
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
	workspacesCommands
}

func newWorkspacesCommands() workspacesCommands {
	return workspacesCommands{
		variables: workspacesVariablesProxyForProduction{},
	}
}

type workspacesVariablesProxy interface {
	list(*tfe.Client, context.Context, string, tfe.VariableListOptions) (*tfe.VariableList, error)
	read(*tfe.Client, context.Context, string, string) (*tfe.Variable, error)
	update(client *tfe.Client, ctx context.Context, workspaceID string, variableID string, opts tfe.VariableUpdateOptions) (*tfe.Variable, error)
}

type workspacesVariablesProxyForProduction struct{}

func (p workspacesVariablesProxyForProduction) list(client *tfe.Client, ctx context.Context, workspaceID string, opts tfe.VariableListOptions) (*tfe.VariableList, error) {
	return client.Variables.List(ctx, workspaceID, opts)
}

func (p workspacesVariablesProxyForProduction) read(client *tfe.Client, ctx context.Context, workspaceID string, variableID string) (*tfe.Variable, error) {
	return client.Variables.Read(ctx, workspaceID, variableID)
}

func (p workspacesVariablesProxyForProduction) update(client *tfe.Client, ctx context.Context, workspaceID string, variableID string, opts tfe.VariableUpdateOptions) (*tfe.Variable, error) {
	return client.Variables.Update(ctx, workspaceID, variableID, opts)
}

type workspacesCommands struct {
	variables workspacesVariablesProxy
}

type dependencyProxies struct {
	os     osProxy
	client clientProxy
}
