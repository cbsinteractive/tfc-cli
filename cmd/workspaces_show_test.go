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

func TestWorkspacesShow(t *testing.T) {
	newDefaultWorkspace := func() *tfe.Workspace {
		return &tfe.Workspace{
			ID:          "some workspace id",
			Description: "some workspace description",
		}
	}
	newDefaultCommandResult := func() string {
		d, _ := json.Marshal(CommandResult{
			Result: WorkspacesShowCommandResult{
				ID:          "some workspace id",
				Description: "some workspace description",
			},
		})
		d = append(d, '\n')
		return string(d)
	}
	testConfigs := []struct {
		description         string
		args                []string
		organization        string
		token               string
		workspace           string
		workspaceShowResult *tfe.Workspace
		expectedError       error
		expectOutput        bool
		expectedOutput      string
	}{
		{
			"show existing workspace",
			[]string{"-workspace", "foo"},
			"some org",
			"some token",
			"foo",
			newDefaultWorkspace(),
			nil,
			true,
			newDefaultCommandResult(),
		},
		{
			"show existing workspace (quiet)",
			[]string{"-workspace", "foo", "-quiet"},
			"some org",
			"some token",
			"foo",
			newDefaultWorkspace(),
			nil,
			false,
			"unused",
		},
		{
			"show missing workspace",
			[]string{"-workspace", "foo"},
			"some org",
			"some token",
			"foo",
			nil,
			errors.New("resource not found"),
			true,
			"resource not found\n",
		},
	}
	for _, d := range testConfigs {
		t.Run(d.description, func(t *testing.T) {
			args := append([]string{"workspaces", "show"}, d.args...)
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
			mockedWorkspacesProxy.On("read", mock.Anything, mock.Anything, d.organization, d.workspace).Return(d.workspaceShowResult, d.expectedError)

			// Code under test
			err := root(
				options,
				args,
				dependencyProxies{
					client: clientProxy{
						workspaces: mockedWorkspacesProxy,
					},
					os: mockedOSProxy,
				},
			)

			// Verify
			if d.expectedError != nil {
				assert.Same(t, d.expectedError, err)
			} else {
				assert.Nil(t, err)
			}
			mockedOSProxy.AssertExpectations(t)
			mockedWorkspacesProxy.AssertExpectations(t)
			if !d.expectOutput {
				assert.Empty(t, buff.String())
			} else {
				assert.Equal(t, d.expectedOutput, buff.String())
			}
		})
	}
}
