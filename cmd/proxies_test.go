package cmd

import (
	"context"
	"errors"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/mock"
)

type mockStateVersionsProxy struct {
	mock.Mock
}

func (m mockStateVersionsProxy) currentWithOptions(client *tfe.Client, ctx context.Context, workspaceID string, options *tfe.StateVersionCurrentOptions) (*tfe.StateVersion, error) {
	args := m.Called(client, ctx, workspaceID, options)
	return args.Get(0).(*tfe.StateVersion), args.Error(1)
}

type mockOSProxy struct {
	mock.Mock
}

func (m mockOSProxy) lookupEnv(key string) (string, bool) {
	args := m.Called(key)
	return args.String(0), args.Bool(1)
}

type mockWorkspacesProxy struct {
	mock.Mock
}

func (m mockWorkspacesProxy) create(client *tfe.Client, ctx context.Context, organization string, options tfe.WorkspaceCreateOptions) (*tfe.Workspace, error) {
	args := m.Called(client, ctx, organization, options)
	return args.Get(0).(*tfe.Workspace), args.Error(1)
}

func (m mockWorkspacesProxy) delete(client *tfe.Client, ctx context.Context, organization string, workspace string) error {
	args := m.Called(client, ctx, organization, workspace)
	return args.Error(0)
}

func (m mockWorkspacesProxy) read(client *tfe.Client, ctx context.Context, organization string, workspace string) (*tfe.Workspace, error) {
	args := m.Called(client, ctx, organization, workspace)
	return args.Get(0).(*tfe.Workspace), args.Error(1)
}

func (m mockWorkspacesProxy) update(client *tfe.Client, ctx context.Context, organization string, workspace string, options tfe.WorkspaceUpdateOptions) (*tfe.Workspace, error) {
	args := m.Called(client, ctx, organization, workspace, options)
	return args.Get(0).(*tfe.Workspace), args.Error(1)
}

func (m mockWorkspacesProxy) removeVCSConnection(client *tfe.Client, ctx context.Context, organization string, workspace string) (*tfe.Workspace, error) {
	args := m.Called(client, ctx, organization, workspace)
	return args.Get(0).(*tfe.Workspace), args.Error(1)
}

type mockWorkspacesVariablesProxy struct {
	mock.Mock
}

func (m mockWorkspacesVariablesProxy) create(client *tfe.Client, ctx context.Context, workspaceID string, options tfe.VariableCreateOptions) (*tfe.Variable, error) {
	args := m.Called(client, ctx, workspaceID, options)
	return args.Get(0).(*tfe.Variable), args.Error(1)
}

func (m mockWorkspacesVariablesProxy) delete(client *tfe.Client, ctx context.Context, workspaceID string, variableID string) error {
	args := m.Called(client, ctx, workspaceID, variableID)
	return args.Error(0)
}

func (m mockWorkspacesVariablesProxy) list(client *tfe.Client, ctx context.Context, workspaceID string, options tfe.VariableListOptions) (*tfe.VariableList, error) {
	args := m.Called(client, ctx, workspaceID, options)
	return args.Get(0).(*tfe.VariableList), args.Error(1)
}

func (m mockWorkspacesVariablesProxy) read(*tfe.Client, context.Context, string, string) (*tfe.Variable, error) {
	return nil, errors.New("not implemented (read)")
}

func (m mockWorkspacesVariablesProxy) update(client *tfe.Client, ctx context.Context, workspaceID string, variableID string, options tfe.VariableUpdateOptions) (*tfe.Variable, error) {
	args := m.Called(client, ctx, workspaceID, variableID, options)
	return args.Get(0).(*tfe.Variable), args.Error(1)
}
