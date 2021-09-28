package cmd

import (
	"bytes"
	"testing"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWorkspacesVariablesDelete(t *testing.T) {
	testConfigs := []struct {
		description  string
		args         []string
		organization string
		token        string
		workspace    string
		workspaceID  string
		variableKey  string
		variableID   string
		deleteError  error
	}{
		{
			"create a new variable",
			[]string{
				"-workspace",
				"foo",
				"-key",
				"bar",
			},
			"some org",
			"some token",
			"foo",
			"some workspace id",
			"bar",
			"some variable id",
			nil,
		},
	}
	for _, d := range testConfigs {
		t.Run(d.description, func(t *testing.T) {
			args := append([]string{"workspaces", "variables", "delete"}, d.args...)
			var buff bytes.Buffer
			options := ExecuteOpts{
				AppName: "tfc-cli",
				Writer:  &buff,
			}
			// Set up expectations
			mockedOSProxy := mockOSProxy{}
			mockedOSProxy.On("lookupEnv", "TFC_ORG").Return(d.organization, true)
			mockedOSProxy.On("lookupEnv", "TFC_TOKEN").Return(d.token, true)
			mockWorkspacesProxy := mockWorkspacesProxy{}
			mockWorkspacesProxy.On("read", mock.Anything, mock.Anything, d.organization, d.workspace).Return(&tfe.Workspace{ID: d.workspaceID}, nil)
			variables := mockWorkspacesVariablesProxy{}
			variables.On("list", mock.Anything, mock.Anything, d.workspaceID, mock.Anything).Return(&tfe.VariableList{
				Items: []*tfe.Variable{
					{
						ID:  d.variableID,
						Key: d.variableKey,
					},
				},
			}, nil)
			variables.On(
				"delete",
				mock.Anything,
				mock.Anything,
				d.workspaceID,
				d.variableID,
			).Return(d.deleteError)

			// Code under test
			err := root(
				options,
				args,
				dependencyProxies{
					client: clientProxy{
						workspaces: mockWorkspacesProxy,
						workspacesCommands: workspacesCommands{
							variables: variables,
						},
					},
					os: mockedOSProxy,
				},
			)

			// Verify
			assert.Nil(t, err)
			mockedOSProxy.AssertExpectations(t)
			mockWorkspacesProxy.AssertExpectations(t)
			variables.AssertExpectations(t)
		})
	}
}
