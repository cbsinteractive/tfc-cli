package cmd

import (
	"context"
	"errors"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/mock"
)

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

// type workspacesProxyForTests struct {
// 	t                *testing.T
// 	organization     string
// 	workspace        string
// 	workspaceID      string
// 	createdWorkspace *tfe.Workspace
// 	createError      error
// }

// func (p workspacesProxyForTests) create(
// 	*tfe.Client,
// 	context.Context,
// 	string,
// 	tfe.WorkspaceCreateOptions,
// ) (*tfe.Workspace, error) {
// 	return p.createdWorkspace, p.createError
// }

// func (p workspacesProxyForTests) delete(
// 	_ *tfe.Client,
// 	_ context.Context,
// 	organization string,
// 	workspace string,
// ) error {
// 	assert.Equal(p.t, p.organization, organization)
// 	assert.Equal(p.t, p.workspace, workspace)
// 	return nil
// }

// func (p workspacesProxyForTests) read(
// 	*tfe.Client,
// 	context.Context,
// 	string,
// 	string,
// ) (*tfe.Workspace, error) {
// 	if p.workspaceID == "" {
// 		return nil, errors.New("resource not found")
// 	}
// 	return &tfe.Workspace{
// 		ID: p.workspaceID,
// 	}, nil
// }

type stateVersionsProxyForTests struct {
	outputs []*tfe.StateVersionOutput
}

func (p stateVersionsProxyForTests) currentWithOptions(
	_ *tfe.Client,
	ctx context.Context,
	workspaceID string,
	options *tfe.StateVersionCurrentOptions,
) (*tfe.StateVersion, error) {
	if p.outputs == nil {
		return nil, errors.New("not implemented")
	}
	return &tfe.StateVersion{
		Outputs: p.outputs,
	}, nil
}

type mockWorkspacesVariablesProxy struct {
	mock.Mock
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
