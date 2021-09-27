package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWorkspacesVariablesList(t *testing.T) {
	testConfigs := []struct {
		description         string
		args                []string
		organization        string
		token               string
		workspaceID         string
		variablesListResult *tfe.VariableList
		variablesListError  error
		expectedResult      WorkspacesVariablesListCommandResult
	}{
		{
			"lists variables for existing workspace",
			[]string{"-workspace", "foo"},
			"some org",
			"some token",
			"some workspace id",
			&tfe.VariableList{
				Items: []*tfe.Variable{
					{
						Key: "foo",
					},
					{
						Key: "bar",
					},
					{
						Key: "baz",
					},
				},
			},
			nil,
			WorkspacesVariablesListCommandResult{
				Result: "foo,bar,baz",
			},
		},
	}
	for _, d := range testConfigs {
		t.Run(d.description, func(t *testing.T) {
			args := append([]string{"workspaces", "variables", "list"}, d.args...)
			var buff bytes.Buffer
			options := ExecuteOpts{
				AppName: "tfc-cli",
				Writer:  &buff,
			}
			mockedOSProxy := mockOSProxy{}
			mockedOSProxy.On("lookupEnv", "TFC_ORG").Return(d.organization, true)
			mockedOSProxy.On("lookupEnv", "TFC_TOKEN").Return(d.token, true)
			mockedVariables := mockWorkspacesVariablesProxy{}
			mockedVariables.On("list", mock.Anything, mock.Anything, d.workspaceID, mock.Anything).Return(d.variablesListResult, d.variablesListError)
			if err := root(
				options,
				args,
				dependencyProxies{
					client: clientProxy{
						workspaces: workspacesProxyForTests{
							workspaceID: d.workspaceID,
						},
						workspacesCommands: workspacesCommands{
							variables: mockedVariables,
						},
					},
					os: mockedOSProxy,
				},
			); err != nil {
				t.Fatal(err)
			}
			// Verify result
			result := WorkspacesVariablesListCommandResult{}
			json.Unmarshal(buff.Bytes(), &result)
			assert.Equal(t, d.expectedResult, result)
		})
	}
}
