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
	return client.StateVersions.ReadCurrentWithOptions(ctx, workspaceID, options)
}

type workspacesProxy interface {
	create(client *tfe.Client, ctx context.Context, organization string, options tfe.WorkspaceCreateOptions) (*tfe.Workspace, error)
	delete(client *tfe.Client, ctx context.Context, organization string, workspace string) error
	read(*tfe.Client, context.Context, string, string) (*tfe.Workspace, error)
	update(client *tfe.Client, ctx context.Context, organization string, workspace string, options tfe.WorkspaceUpdateOptions) (*tfe.Workspace, error)
	removeVCSConnection(client *tfe.Client, ctx context.Context, organization string, workspace string) (*tfe.Workspace, error)
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

func (p workspacesProxyForProduction) update(client *tfe.Client, ctx context.Context, organization string, workspace string, options tfe.WorkspaceUpdateOptions) (*tfe.Workspace, error) {
	return client.Workspaces.Update(ctx, organization, workspace, options)
}

func (p workspacesProxyForProduction) removeVCSConnection(client *tfe.Client, ctx context.Context, organization string, workspace string) (*tfe.Workspace, error) {
	return client.Workspaces.RemoveVCSConnection(ctx, organization, workspace)
}

type clientProxy struct {
	stateVersions stateVersionsProxy
	workspaces    workspacesProxy
	workspacesCommands
}

func newWorkspacesCommands() workspacesCommands {
	return workspacesCommands{
		tags:      workspacesTagsProxyForProduction{},
		variables: workspacesVariablesProxyForProduction{},
	}
}

type workspacesVariablesProxy interface {
	create(client *tfe.Client, ctx context.Context, workspaceID string, options tfe.VariableCreateOptions) (*tfe.Variable, error)
	delete(client *tfe.Client, ctx context.Context, workspaceID string, variableID string) error
	list(client *tfe.Client, ctx context.Context, workspaceID string, options *tfe.VariableListOptions) (*tfe.VariableList, error)
	read(client *tfe.Client, ctx context.Context, workspaceID string, variableID string) (*tfe.Variable, error)
	update(client *tfe.Client, ctx context.Context, workspaceID string, variableID string, options tfe.VariableUpdateOptions) (*tfe.Variable, error)
}

type workspacesVariablesProxyForProduction struct{}

func (p workspacesVariablesProxyForProduction) create(client *tfe.Client, ctx context.Context, workspaceID string, options tfe.VariableCreateOptions) (*tfe.Variable, error) {
	return client.Variables.Create(ctx, workspaceID, options)
}

func (p workspacesVariablesProxyForProduction) delete(client *tfe.Client, ctx context.Context, workspaceID string, variableID string) error {
	return client.Variables.Delete(ctx, workspaceID, variableID)
}

func (p workspacesVariablesProxyForProduction) list(client *tfe.Client, ctx context.Context, workspaceID string, opts *tfe.VariableListOptions) (*tfe.VariableList, error) {
	return client.Variables.List(ctx, workspaceID, opts)
}

func (p workspacesVariablesProxyForProduction) read(client *tfe.Client, ctx context.Context, workspaceID string, variableID string) (*tfe.Variable, error) {
	return client.Variables.Read(ctx, workspaceID, variableID)
}

func (p workspacesVariablesProxyForProduction) update(client *tfe.Client, ctx context.Context, workspaceID string, variableID string, opts tfe.VariableUpdateOptions) (*tfe.Variable, error) {
	return client.Variables.Update(ctx, workspaceID, variableID, opts)
}

type workspacesTagsProxy interface {
	create(client *tfe.Client, ctx context.Context, workspaceID string, options tfe.WorkspaceAddTagsOptions) error
	delete(client *tfe.Client, ctx context.Context, workspaceID string, options tfe.WorkspaceRemoveTagsOptions) error
}

type workspacesTagsProxyForProduction struct{}

func (p workspacesTagsProxyForProduction) create(client *tfe.Client, ctx context.Context, workspaceID string, options tfe.WorkspaceAddTagsOptions) error {
	return client.Workspaces.AddTags(ctx, workspaceID, options)
}

func (p workspacesTagsProxyForProduction) delete(client *tfe.Client, ctx context.Context, workspaceID string, options tfe.WorkspaceRemoveTagsOptions) error {
	return client.Workspaces.RemoveTags(ctx, workspaceID, options)
}

type workspacesCommands struct {
	tags      workspacesTagsProxy
	variables workspacesVariablesProxy
}

type dependencyProxies struct {
	os     osProxy
	client clientProxy
}
