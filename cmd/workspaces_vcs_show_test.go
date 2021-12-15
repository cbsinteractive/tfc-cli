package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWorkspacesVCSShow(t *testing.T) {
	testConfigs := []struct {
		description          string
		args                 []string
		organization         string
		token                string
		workspace            string
		workspaceID          string
		workspaceReadResult  *tfe.Workspace
		workspaceReadError   error
		expectedResultObject CommandResult
	}{
		{
			"VCS repo object is nil",
			[]string{"-workspace", "some workspace"},
			"some org",
			"some token",
			"some workspace",
			"some workspace id",
			&tfe.Workspace{
				ID: "some workspace id",
			},
			nil,
			CommandResult{
				Result: "VCS repo not set",
			},
		},
	}
	for _, d := range testConfigs {
		t.Run(d.description, func(t *testing.T) {
			args := append([]string{"workspaces", "vcs", "show"}, d.args...)
			var buff bytes.Buffer
			options := ExecuteOpts{
				AppName: "tfc-cli",
				Stdout:  &buff,
			}
			mockedOSProxy := mockOSProxy{}
			mockedOSProxy.On("lookupEnv", "TFC_ORG").Return(d.organization, true)
			mockedOSProxy.On("lookupEnv", "TFC_TOKEN").Return(d.token, true)
			mockedWorkspacesProxy := mockWorkspacesProxy{}
			mockedWorkspacesProxy.On("read", mock.Anything, mock.Anything, d.organization, d.workspace).Return(d.workspaceReadResult, d.workspaceReadError)

			// Code under test
			err := root(
				options,
				args,
				dependencyProxies{
					client: clientProxy{
						workspaces:         mockedWorkspacesProxy,
						workspacesCommands: workspacesCommands{
							// variables: mockedVariablesProxy,
						},
					},
					os: mockedOSProxy,
				},
			)

			// Verify
			assert.Nil(t, err)
			r := CommandResult{}
			json.Unmarshal(buff.Bytes(), &r)
			assert.Equal(t, d.expectedResultObject, r)
		})
	}
}
