package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWorkspacesUpdate(t *testing.T) {
	newWorkspaceUpdateOptions := (func(description string) tfe.WorkspaceUpdateOptions {
		return tfe.WorkspaceUpdateOptions{
			Description: &description,
		}
	})
	testConfigs := []struct {
		description           string
		args                  []string
		organization          string
		token                 string
		workspace             string
		workspaceDescription  string
		workspaceReadResult   *tfe.Workspace
		workspaceReadError    error
		workspaceUpdateResult *tfe.Workspace
		workspaceUpdateError  error
		expectedResultObject  *CommandResult
		expectedError         error
	}{
		{
			"Updates description",
			[]string{"-workspace", "someworkspace", "-description", "new description"},
			"someorg",
			"sometoken",
			"someworkspace",
			"new description",
			&tfe.Workspace{
				Description: "old description",
			},
			nil,
			&tfe.Workspace{
				ID:          "foo",
				Description: "new description",
			},
			nil,
			&CommandResult{
				Result: WorkspacesUpdateCommandResult{
					ID:          "foo",
					Description: "new description",
				},
			},
			nil,
		},
	}
	for _, test := range testConfigs {
		t.Run(test.description, func(t *testing.T) {
			args := append([]string{"workspaces", "update"}, test.args...)
			var stdBuff, errBuff bytes.Buffer
			options := ExecuteOpts{
				AppName: "tfc-cli",
				Stdout:  &stdBuff,
				Stderr:  &errBuff,
			}
			// Set up expectations
			mockedOSProxy := mockOSProxy{}
			mockedOSProxy.On("lookupEnv", "TFC_ORG").Return(test.organization, true)
			mockedOSProxy.On("lookupEnv", "TFC_TOKEN").Return(test.token, true)
			mockedWorkspacesProxy := mockWorkspacesProxy{}
			mockedWorkspacesProxy.On(
				"read",
				mock.Anything,
				mock.Anything,
				test.organization,
				test.workspace,
			).Return(
				test.workspaceReadResult,
				test.workspaceReadError,
			)
			mockedWorkspacesProxy.On(
				"update",
				mock.Anything,
				mock.Anything,
				test.organization,
				test.workspace,
				newWorkspaceUpdateOptions(
					test.workspaceDescription,
				),
			).Return(
				test.workspaceUpdateResult,
				test.workspaceUpdateError,
			)

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
			mockedOSProxy.AssertExpectations(t)
			mockedWorkspacesProxy.AssertExpectations(t)

			if test.expectedError == nil {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expectedError.Error())
			}
			if test.expectedResultObject != nil {
				expectedOutput, _ := json.Marshal(test.expectedResultObject)
				expectedOutput = append(expectedOutput, '\n')
				assert.Equal(t, string(expectedOutput), stdBuff.String())
			} else {
				assert.Empty(t, stdBuff.String())
			}
		})
	}
}
