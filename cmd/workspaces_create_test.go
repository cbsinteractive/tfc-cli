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
	newWorkspaceCreateOptions := (func(name string, workingDirectory string, vcsRepoOpts *tfe.VCSRepoOptions) tfe.WorkspaceCreateOptions {
		description := "Created by tfc-cli"
		opts := tfe.WorkspaceCreateOptions{
			Name:             &name,
			Description:      &description,
			WorkingDirectory: &workingDirectory,
			VCSRepo:          vcsRepoOpts,
		}
		return opts
	})
	newVCSRepoOptions := (func(identifier string, branch string, oauthTokenId string) *tfe.VCSRepoOptions {
		return &tfe.VCSRepoOptions{
			Identifier:   &identifier,
			Branch:       &branch,
			OAuthTokenID: &oauthTokenId,
		}
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
		workingDirectory      string
		expectedVCSRepoOpts   *tfe.VCSRepoOptions
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
			"",
			nil,
			&tfe.Workspace{
				ID:          "someid",
				Description: "some description",
			},
			nil,
			newDefaultCommandResult(),
			nil,
		},
		{
			"with working directory",
			[]string{"-workspace", "foo", "-working-directory", "somedirectory"},
			"some org",
			"some token",
			"foo",
			"somedirectory",
			nil,
			&tfe.Workspace{
				ID:          "someid",
				Description: "some description",
			},
			nil,
			newDefaultCommandResult(),
			nil,
		},
		{
			"with VCS options",
			[]string{
				"-workspace", "foo", "-vcs-identifier", "someorg/somerepo",
				"-vcs-branch", "somebranch", "-vcs-oauth-token-id", "someoauthtokenid",
			},
			"some org",
			"some token",
			"foo",
			"",
			newVCSRepoOptions("someorg/somerepo", "somebranch", "someoauthtokenid"),
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
				newWorkspaceCreateOptions(
					d.workspace,
					d.workingDirectory,
					d.expectedVCSRepoOpts,
				)).Return(
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
