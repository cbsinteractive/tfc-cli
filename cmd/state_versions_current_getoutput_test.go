package cmd

import (
	"bytes"
	"testing"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStateVersionsCurrentGetOutput(t *testing.T) {
	testConfigs := []struct {
		name                string
		organization        string
		token               string
		workspace           string
		workspaceId         string
		workspaceReadResult *tfe.Workspace
		workspaceReadError  error
		outputName          string
		outputs             []*tfe.StateVersionOutput
		expectedValue       string
	}{
		{
			"output variable found",
			"some org",
			"some token",
			"some workspace",
			"some workspace id",
			&tfe.Workspace{},
			nil,
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
			// Set up expectations
			mockedOSProxy := mockOSProxy{}
			mockedOSProxy.On("lookupEnv", "TFC_ORG").Return(d.organization, true)
			mockedOSProxy.On("lookupEnv", "TFC_TOKEN").Return(d.token, true)
			mockedWorkspacesProxy := mockWorkspacesProxy{}
			mockedWorkspacesProxy.On("read", mock.Anything, mock.Anything, d.organization, d.workspace).Return(d.workspaceReadResult, d.workspaceReadError)
			if err := root(
				options,
				args,
				dependencyProxies{
					client: clientProxy{
						stateVersions: stateVersionsProxyForTests{
							outputs: d.outputs,
						},
						workspaces: mockedWorkspacesProxy,
					},
					os: mockedOSProxy,
				},
			); err != nil {
				t.Fatal(err)
			}
			assert.Contains(t, buff.String(), d.expectedValue)
		})
	}
}
