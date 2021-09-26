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

func (p osProxyForTests) lookupEnv(key string) (string, bool) {
	if _, ok := p.envVars[key]; !ok {
		return "", false
	}
	return p.envVars[key], true
}

type workspacesProxyForTests struct {
	t                *testing.T
	organization     string
	workspace        string
	workspaceID      string
	createdWorkspace *tfe.Workspace
	createError      error
}

func (p workspacesProxyForTests) create(
	*tfe.Client,
	context.Context,
	string,
	tfe.WorkspaceCreateOptions,
) (*tfe.Workspace, error) {
	return p.createdWorkspace, p.createError
}

func (p workspacesProxyForTests) delete(
	_ *tfe.Client,
	_ context.Context,
	organization string,
	workspace string,
) error {
	assert.Equal(p.t, p.organization, organization)
	assert.Equal(p.t, p.workspace, workspace)
	return nil
}

func (p workspacesProxyForTests) read(
	*tfe.Client,
	context.Context,
	string,
	string,
) (*tfe.Workspace, error) {
	if p.workspaceID == "" {
		return nil, errors.New("resource not found")
	}
	return &tfe.Workspace{
		ID: p.workspaceID,
	}, nil
}

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

func newDefaultEnvForTests() map[string]string {
	return map[string]string{
		"TFC_TOKEN": "some token",
		"TFC_ORG":   "some org",
	}
}

type workspacesVariablesProxyForTesting struct {
	t                    *testing.T
	listVariables        *tfe.VariableList
	listError            error
	readVariable         *tfe.Variable
	readError            error
	updateWorkspaceID    string
	updateVariableID     string
	updateResultVariable *tfe.Variable
	updateError          error
}

func newWorkspacesVariablesProxyForTesting(t *testing.T) workspacesVariablesProxyForTesting {
	return workspacesVariablesProxyForTesting{t: t}
}

func (p workspacesVariablesProxyForTesting) list(*tfe.Client, context.Context, string, tfe.VariableListOptions) (*tfe.VariableList, error) {
	return p.listVariables, p.listError
}

func (p workspacesVariablesProxyForTesting) read(client *tfe.Client, ctx context.Context, workspaceID string, variableID string) (*tfe.Variable, error) {
	return p.readVariable, p.readError
}

func (p workspacesVariablesProxyForTesting) update(client *tfe.Client, ctx context.Context, workspaceID string, variableID string, opts tfe.VariableUpdateOptions) (*tfe.Variable, error) {
	assert.Equal(p.t, p.updateWorkspaceID, workspaceID)
	assert.Equal(p.t, p.updateVariableID, variableID)
	return p.updateResultVariable, p.updateError
}
