package cmd

import (
	"bytes"
	"testing"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
)

func TestStateVersionsCurrentGetOutput(t *testing.T) {
	testConfigs := []struct {
		name          string
		env           map[string]string
		workspaceId   string
		outputName    string
		outputs       []*tfe.StateVersionOutput
		expectedValue string
	}{
		{
			"output variable found",
			newDefaultEnvForTests(),
			"some workspace id",
			"foo",
			[]*tfe.StateVersionOutput{
				{
					Name:  "foo",
					Value: "some value",
				},
			},
			"some value"},
	}
	for _, d := range testConfigs {
		t.Run(d.name, func(t *testing.T) {
			args := []string{
				"stateversions",
				"current",
				"getoutput",
				"-workspace",
				"some workspace",
				"-name",
				d.outputName,
			}
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
						stateVersions: stateVersionsProxyForTests{
							outputs: d.outputs,
						},
						workspaces: workspacesProxyForTests{
							workspaceId: d.workspaceId,
						},
					},
					os: osProxyForTests{
						envVars: d.env,
					},
				},
			); err != nil {
				t.Fatal(err)
			}
			assert.Contains(t, buff.String(), d.expectedValue)
		})
	}
}
