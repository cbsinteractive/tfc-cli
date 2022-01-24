package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWorkspacesCreate(t *testing.T) {
	newWorkspaceCreateOptions := (func(name string) tfe.WorkspaceCreateOptions {
		description := "Created by tfc-cli"
		opts := tfe.WorkspaceCreateOptions{
			Name:        &name,
			Description: &description,
		}
		return opts
	})
	newDefaultCommandResult := (func() *CommandResult {
		return &CommandResult{
			Result: WorkspacesCreateCommandResult{
				ID:          "someid",
				Description: "some description",
			},
		}
	})
	testConfigs := []struct {
		description           string
		args                  []string
		organization          string
		token                 string
		workspace             string
		workspaceCreateResult *tfe.Workspace
		workspaceCreateError  error
		expectedResultObject  *CommandResult
		expectedError         error
	}{
		{
			"workspace created",
			[]string{"-workspace", "foo"},
			"some org",
			"some token",
			"foo",
			&tfe.Workspace{
				ID:          "someid",
				Description: "some description",
			},
			nil,
			newDefaultCommandResult(),
			nil,
		},
	}
	for _, d := range testConfigs {
		t.Run(d.description, func(t *testing.T) {
			args := append([]string{"workspaces", "create"}, d.args...)
			var stdBuff, errBuff bytes.Buffer
			options := ExecuteOpts{
				AppName: "tfc-cli",
				Stdout:  &stdBuff,
				Stderr:  &errBuff,
			}
			// Set up expectations
			mockedOSProxy := mockOSProxy{}
			mockedOSProxy.On("lookupEnv", "TFC_ORG").Return(d.organization, true)
			mockedOSProxy.On("lookupEnv", "TFC_TOKEN").Return(d.token, true)
			mockedWorkspacesProxy := mockWorkspacesProxy{}
			mockedWorkspacesProxy.On(
				"create",
				mock.Anything,
				mock.Anything,
				d.organization,
				newWorkspaceCreateOptions(d.workspace),
			).Return(
				d.workspaceCreateResult,
				d.workspaceCreateError,
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

			if d.expectedError == nil {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, d.expectedError.Error())
			}
			if d.expectedResultObject != nil {
				expectedOutput, _ := json.Marshal(d.expectedResultObject)
				expectedOutput = append(expectedOutput, '\n')
				assert.Equal(t, string(expectedOutput), stdBuff.String())
			} else {
				assert.Empty(t, stdBuff.String())
			}
		})
	}
}
