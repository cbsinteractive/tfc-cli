package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWorkspacesVariablesCreate(t *testing.T) {
	testConfigs := []struct {
		description         string
		args                []string
		organization        string
		token               string
		workspace           string
		workspaceID         string
		variableKey         string
		variableID          string
		expectedValue       string
		expectedDescription string
		expectedCategory    tfe.CategoryType
		expectedSensitive   bool
		expectedHCL         bool
	}{
		{
			"create a new variable",
			[]string{
				"-workspace",
				"foo",
				"-key",
				"bar",
				"-value",
				"baz",
				"-description",
				"quux",
				"-category",
				"terraform",
				"-sensitive=false",
				"-hcl=false",
			},
			"some org",
			"some token",
			"foo",
			"some workspace id",
			"bar",
			"some variable id",
			"baz",
			"quux",
			"terraform",
			false,
			false,
		},
	}
	for _, d := range testConfigs {
		t.Run(d.description, func(t *testing.T) {
			args := append([]string{"workspaces", "variables", "create"}, d.args...)
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
			variables.On(
				"create",
				mock.Anything,
				mock.Anything,
				d.workspaceID,
				tfe.VariableCreateOptions{
					Key:         &d.variableKey,
					Value:       &d.expectedValue,
					Description: &d.expectedDescription,
					Category:    &d.expectedCategory,
					Sensitive:   &d.expectedSensitive,
					HCL:         &d.expectedHCL,
				}).Return(&tfe.Variable{}, nil)

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
			result := WorkspacesVariablesUpdateValueCommandResult{}
			assert.Nil(t, json.Unmarshal(buff.Bytes(), &result))
		})
	}
}
