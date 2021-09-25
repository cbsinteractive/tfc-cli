package cmd

import (
	"context"
	"errors"
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
	read(*tfe.Client, context.Context, string, string) (*tfe.Workspace, error)
}

type workspacesProxyForProduction struct{}

func newWorkspacesProxy() workspacesProxyForProduction {
	return workspacesProxyForProduction{}
}

func (p workspacesProxyForProduction) create(client *tfe.Client, ctx context.Context, organization string, opts tfe.WorkspaceCreateOptions) (*tfe.Workspace, error) {
	return client.Workspaces.Create(ctx, organization, opts)
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

type osProxyForTests struct {
	envVars map[string]string
}

func (c osProxyForTests) lookupEnv(key string) (string, bool) {
	if _, ok := c.envVars[key]; !ok {
		return "", false
	}
	return c.envVars[key], true
}

type workspacesProxyForTests struct {
	workspaceId      string
	createdWorkspace *tfe.Workspace
	createError      error
}

func (c workspacesProxyForTests) create(
	*tfe.Client,
	context.Context,
	string,
	tfe.WorkspaceCreateOptions,
) (*tfe.Workspace, error) {
	return c.createdWorkspace, c.createError
}

func (c workspacesProxyForTests) read(
	*tfe.Client,
	context.Context,
	string,
	string,
) (*tfe.Workspace, error) {
	if c.workspaceId == "" {
		return nil, errors.New("resource not found")
	}
	return &tfe.Workspace{
		ID: c.workspaceId,
	}, nil
}

type stateVersionsProxyForTests struct {
	outputs []*tfe.StateVersionOutput
}

func (c stateVersionsProxyForTests) currentWithOptions(
	_ *tfe.Client,
	ctx context.Context,
	workspaceID string,
	options *tfe.StateVersionCurrentOptions,
) (*tfe.StateVersion, error) {
	if c.outputs == nil {
		return nil, errors.New("not implemented")
	}
	return &tfe.StateVersion{
		Outputs: c.outputs,
	}, nil
}

func newDefaultEnvForTests() map[string]string {
	return map[string]string{
		"TFC_TOKEN": "some token",
		"TFC_ORG":   "some org",
	}
}
