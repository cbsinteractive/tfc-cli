package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
)

func TestWorkspacesVariablesUpdate(t *testing.T) {
	testConfigs := []struct {
		description       string
		args              []string
		workspaceID       string
		newVariablesProxy func(*testing.T) workspacesVariablesProxy
		expectedResult    WorkspacesVariablesUpdateCommandResult
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
				"-description",
				"some description",
			},
			"some workspace id",
			(func(t *testing.T) workspacesVariablesProxy {
				p := newWorkspacesVariablesProxyForTesting(t)
				p.listVariables = &tfe.VariableList{
					Items: []*tfe.Variable{
						{
							ID:  "some variable id",
							Key: "bar",
						},
					},
				}
				p.updateWorkspaceID = "some workspace id"
				p.updateVariableID = "some variable id"
				expectedValue := "baz"
				p.expectedVariableUpdateOptions = tfe.VariableUpdateOptions{
					Value: &expectedValue,
				}
				p.updateResultVariable = &tfe.Variable{
					ID:          "some variable id",
					Key:         "bar",
					Value:       "baz",
					Description: "some description",
					Category:    "terraform",
					HCL:         false,
					Sensitive:   true,
				}
				return p
			}),
			WorkspacesVariablesUpdateCommandResult{
				Result: &tfe.Variable{
					ID:          "some variable id",
					Key:         "bar",
					Value:       "baz",
					Description: "some description",
					Category:    "terraform",
					Sensitive:   true,
					HCL:         false,
				},
			},
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
			if err := root(
				options,
				args,
				dependencyProxies{
					client: clientProxy{
						workspaces: workspacesProxyForTests{
							workspaceID: d.workspaceID,
						},
						workspacesCommands: workspacesCommands{
							variables: d.newVariablesProxy(t),
						},
					},
					os: osProxyForTests{
						envVars: newDefaultEnvForTests(),
					},
				},
			); err != nil {
				t.Fatal(err)
			}
			// Verify result
			result := WorkspacesVariablesUpdateCommandResult{}
			json.Unmarshal(buff.Bytes(), &result)
			assert.Equal(t, d.expectedResult, result)
		})
	}
}