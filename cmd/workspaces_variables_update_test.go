package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	// "github.com/stretchr/testify/mock"
)

func TestWorkspacesVariablesUpdate(t *testing.T) {
	testConfigs := []struct {
		description string
		args        []string
		workspaceID string
	}{
		{
			"update existing variable",
			[]string{
				"-workspace",
				"foo",
				"-key",
				"bar",
				"-value",
				"baz",
				"-category",
				"terraform",
				"-sensitive=true",
				"-hcl=false",
				"-description=\"some description\"",
			},
			"some workspace id",
		},
	}
	for _, d := range testConfigs {
		t.Run(d.description, func(t *testing.T) {
			args := append([]string{"workspaces", "variables", "update"}, d.args...)
			var buff bytes.Buffer
			options := ExecuteOpts{
				AppName: "tfc-cli",
				Writer:  &buff,
			}
			// Set up expectations
			mockOSProxy := mockOSProxy{}
			mockOSProxy.On("lookupEnv", "TFC_ORG").Return("some org", true)
			mockOSProxy.On("lookupEnv", "TFC_TOKEN").Return("some token", true)
			mockWorkspacesProxy := mockWorkspacesProxy{}
			mockWorkspacesProxy.On("read", mock.Anything, mock.Anything, "some org", "foo").Return(&tfe.Workspace{ID: "some workspace id"}, nil)
			variables := mockWorkspacesVariablesProxy{}
			variables.On("list", mock.Anything, mock.Anything, "some workspace id", mock.Anything).Return(&tfe.VariableList{
				Items: []*tfe.Variable{
					{
						ID:  "some variable id",
						Key: "bar",
					},
				},
			}, nil)
			expectedValue := "baz"
			expectedDescription := "\"some description\""
			variables.On("update", mock.Anything, mock.Anything, "some workspace id", "some variable id", tfe.VariableUpdateOptions{Value: &expectedValue, Description: &expectedDescription}).Return(&tfe.Variable{}, nil)

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
					os: mockOSProxy,
				},
			)

			// Verify
			assert.Nil(t, err)
			mockOSProxy.AssertExpectations(t)
			mockWorkspacesProxy.AssertExpectations(t)
			variables.AssertExpectations(t)
			result := WorkspacesVariablesUpdateCommandResult{}
			assert.Nil(t, json.Unmarshal(buff.Bytes(), &result))
		})
	}
}
