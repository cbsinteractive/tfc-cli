package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
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
		expectedResultObject *CommandResult
		expectedError        error
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
			&CommandResult{
				Result: "VCS repo not set",
			},
			nil,
		},
		{
			"invalid workspace",
			[]string{"-workspace", "some workspace"},
			"some org",
			"some token",
			"some workspace",
			"some workspace id",
			nil,
			errors.New("foo"),
			nil,
			errors.New("foo"),
		},
	}
	for _, d := range testConfigs {
		t.Run(d.description, func(t *testing.T) {
			args := append([]string{"workspaces", "vcs", "show"}, d.args...)
			var stdBuff, errBuff bytes.Buffer
			options := ExecuteOpts{
				AppName: "tfc-cli",
				Stdout:  &stdBuff,
				Stderr:  &errBuff,
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
			if d.expectedError == nil {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, d.expectedError.Error())
			}
			if d.expectedResultObject != nil {
				r := CommandResult{}
				json.Unmarshal(stdBuff.Bytes(), &r)
				assert.Equal(t, *d.expectedResultObject, r)
			} else {
				assert.Empty(t, stdBuff.String())
			}
		})
	}
}
