package cmd

import (
	"bytes"
	"testing"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
)

func TestWorkspacesVariablesList(t *testing.T) {
	testConfigs := []struct {
		description   string
		args          []string
		workspaceID   string
		listVariables *tfe.VariableList
		listError     error
		expectedList  string
	}{
		{
			"lists variables for existing workspace",
			[]string{"-workspace", "foo"},
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
			"foo,bar,baz",
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
			if err := root(
				options,
				args,
				dependencyProxies{
					client: clientProxy{
						workspaces: workspacesProxyForTests{
							workspaceId: d.workspaceID,
						},
						workspacesCommands: workspacesCommands{
							variables: newWorkspacesVariablesProxyForTesting(d.listVariables, d.listError),
						},
					},
					os: osProxyForTests{
						envVars: newDefaultEnvForTests(),
					},
				},
			); err != nil {
				t.Fatal(err)
			}
			assert.Contains(t, buff.String(), d.expectedList)
		})
	}
}
