package cmd

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
)

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
	t                *testing.T
	organization     string
	workspace        string
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

func (c workspacesProxyForTests) delete(
	_ *tfe.Client,
	_ context.Context,
	organization string,
	workspace string,
) error {
	assert.Equal(c.t, c.organization, organization)
	assert.Equal(c.t, c.workspace, workspace)
	return nil
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

type workspacesVariablesProxyForTesting struct {
	listVariables *tfe.VariableList
	listError     error
}

func newWorkspacesVariablesProxyForTesting(listVariables *tfe.VariableList, listError error) workspacesVariablesProxyForTesting {
	return workspacesVariablesProxyForTesting{
		listVariables: listVariables,
		listError:     listError,
	}
}

func (p workspacesVariablesProxyForTesting) list(*tfe.Client, context.Context, string, tfe.VariableListOptions) (*tfe.VariableList, error) {
	return p.listVariables, p.listError
}
