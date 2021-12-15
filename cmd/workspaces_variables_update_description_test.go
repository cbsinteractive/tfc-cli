package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWorkspacesVariablesUpdateDescription(t *testing.T) {
	testConfigs := []struct {
		description         string
		args                []string
		organization        string
		token               string
		workspace           string
		workspaceID         string
		variableKey         string
		variableID          string
		expectedDescription string
		updateResult        *tfe.Variable
		updateError         error
		expectedOutput      CommandResult
	}{
		{
			"update existing variable value",
			[]string{
				"-workspace",
				"foo",
				"-key",
				"bar",
				"-description",
				"baz",
			},
			"some org",
			"some token",
			"foo",
			"some workspace id",
			"bar",
			"some variable id",
			"baz",
			&tfe.Variable{
				ID:          "some variable id",
				Key:         "bar",
				Description: "baz",
			},
			nil,
			CommandResult{
				Result: map[string]interface{}{
					"id":          "some variable id",
					"key":         "bar",
					"description": "baz",
				},
			},
		},
	}
	for _, d := range testConfigs {
		t.Run(d.description, func(t *testing.T) {
			args := append([]string{"workspaces", "variables", "update", "description"}, d.args...)
			var buff bytes.Buffer
			options := ExecuteOpts{
				AppName: "tfc-cli",
				Stdout:  &buff,
			}
			// Set up expectations
			mockedOSProxy := mockOSProxy{}
			mockedOSProxy.On("lookupEnv", "TFC_ORG").Return(d.organization, true)
			mockedOSProxy.On("lookupEnv", "TFC_TOKEN").Return(d.token, true)
			mockWorkspacesProxy := mockWorkspacesProxy{}
			mockWorkspacesProxy.On("read", mock.Anything, mock.Anything, d.organization, d.workspace).Return(&tfe.Workspace{ID: d.workspaceID}, nil)
			mockedVariablesProxy := mockWorkspacesVariablesProxy{}
			mockedVariablesProxy.On("list", mock.Anything, mock.Anything, d.workspaceID, mock.Anything).Return(&tfe.VariableList{
				Items: []*tfe.Variable{
					{
						ID:  d.variableID,
						Key: d.variableKey,
					},
				},
			}, nil)
			mockedVariablesProxy.On(
				"update",
				mock.Anything,
				mock.Anything,
				d.workspaceID,
				d.variableID,
				tfe.VariableUpdateOptions{
					Description: &d.expectedDescription,
				}).Return(d.updateResult, d.updateError)

			// Code under test
			err := root(
				options,
				args,
				dependencyProxies{
					client: clientProxy{
						workspaces: mockWorkspacesProxy,
						workspacesCommands: workspacesCommands{
							variables: mockedVariablesProxy,
						},
					},
					os: mockedOSProxy,
				},
			)

			// Verify
			assert.Nil(t, err)
			mockedOSProxy.AssertExpectations(t)
			mockWorkspacesProxy.AssertExpectations(t)
			mockedVariablesProxy.AssertExpectations(t)
			var result CommandResult
			err = json.Unmarshal(buff.Bytes(), &result)
			assert.Nil(t, err)
			assert.Equal(t, d.expectedOutput, result)
		})
	}
}
