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
		name                     string
		organization             string
		token                    string
		workspace                string
		workspaceID              string
		workspaceReadResult      *tfe.Workspace
		workspaceReadError       error
		outputName               string
		currentWithOptionsResult *tfe.StateVersion
		currentWithOptionsError  error
		expectedValue            string
	}{
		{
			"output variable found",
			"some org",
			"some token",
			"some workspace",
			"some workspace id",
			&tfe.Workspace{
				ID: "some workspace id",
			},
			nil,
			"foo",
			&tfe.StateVersion{
				Outputs: []*tfe.StateVersionOutput{
					{
						Name:  "foo",
						Value: "some value",
					},
				},
			},
			nil,
			"some value",
		},
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
				Stdout:  &buff,
			}
			// Set up expectations
			mockedOSProxy := mockOSProxy{}
			mockedOSProxy.On("lookupEnv", "TFC_ORG").Return(d.organization, true)
			mockedOSProxy.On("lookupEnv", "TFC_TOKEN").Return(d.token, true)
			mockedWorkspacesProxy := mockWorkspacesProxy{}
			mockedWorkspacesProxy.On("read", mock.Anything, mock.Anything, d.organization, d.workspace).Return(d.workspaceReadResult, d.workspaceReadError)
			mockedStateVersionsProxy := mockStateVersionsProxy{}
			mockedStateVersionsProxy.On("currentWithOptions", mock.Anything, mock.Anything, d.workspaceID, &tfe.StateVersionCurrentOptions{Include: "outputs"}).Return(d.currentWithOptionsResult, d.currentWithOptionsError)
			// Code under test
			err := root(
				options,
				args,
				dependencyProxies{
					client: clientProxy{
						stateVersions: mockedStateVersionsProxy,
						workspaces:    mockedWorkspacesProxy,
					},
					os: mockedOSProxy,
				},
			)

			// Verify
			assert.Nil(t, err)
			assert.Contains(t, buff.String(), d.expectedValue)
		})
	}
}
