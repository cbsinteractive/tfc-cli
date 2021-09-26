package cmd

import (
	"bytes"
	"testing"

	"github.com/hashicorp/go-tfe"
)

func TestWorkspacesVariablesUpdate(t *testing.T) {
	testConfigs := []struct {
		description               string
		args                      []string
		workspaceID               string
		listVariables             *tfe.VariableList
		updatedVariable           *tfe.Variable
		expectedUpdateWorkspaceID string
		expectedUpdateVariableID  string
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
				"-sensitive",
				"true",
				"-hcl",
				"false",
			},
			"some workspace id",
			&tfe.VariableList{
				Items: []*tfe.Variable{
					{
						ID:  "some variable id",
						Key: "bar",
					},
				},
			},
			&tfe.Variable{},
			"some workspace id",
			"some variable id",
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
			variablesProxy := newWorkspacesVariablesProxyForTesting(t)
			variablesProxy.listVariables = d.listVariables
			variablesProxy.updateResultVariable = d.updatedVariable
			variablesProxy.updateWorkspaceID = d.expectedUpdateWorkspaceID
			variablesProxy.updateVariableID = d.expectedUpdateVariableID
			if err := root(
				options,
				args,
				dependencyProxies{
					client: clientProxy{
						workspaces: workspacesProxyForTests{
							workspaceID: d.workspaceID,
						},
						workspacesCommands: workspacesCommands{
							variables: variablesProxy,
						},
					},
					os: osProxyForTests{
						envVars: newDefaultEnvForTests(),
					},
				},
			); err != nil {
				t.Fatal(err)
			}
		})
	}
}
