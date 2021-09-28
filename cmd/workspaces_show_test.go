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
	newDefaultCommandResult := func() CommandResult {
		return CommandResult{
			Result: map[string]interface{}{
				"id":          "some workspace id",
				"description": "some workspace description",
			},
		}
	}
	testConfigs := []struct {
		description         string
		args                []string
		organization        string
		token               string
		workspace           string
		workspaceShowResult *tfe.Workspace
		workspaceShowError  error
		expectOutput        bool
		expectedOutput      CommandResult
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
			newDefaultCommandResult(),
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
			CommandResult{
				Error: "resource not found",
			},
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
			mockedWorkspacesProxy.On("read", mock.Anything, mock.Anything, d.organization, d.workspace).Return(d.workspaceShowResult, d.workspaceShowError)

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
			if d.workspaceShowError != nil {
				assert.Same(t, d.workspaceShowError, err)
			} else {
				assert.Nil(t, err)
			}
			mockedOSProxy.AssertExpectations(t)
			mockedWorkspacesProxy.AssertExpectations(t)
			if !d.expectOutput {
				assert.Empty(t, buff.String())
			} else {
				var result CommandResult
				err := json.Unmarshal(buff.Bytes(), &result)
				assert.Nil(t, err)
				assert.Equal(t, d.expectedOutput, result)
			}
		})
	}
}
