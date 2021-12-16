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

func TestWorkspacesUpdate(t *testing.T) {
	type resultConfig struct {
		resultObj *CommandResult
		resultErr error
	}

	type updateConfig struct {
		resultObj *tfe.Workspace
		resultErr error
	}

	type vcsConfig struct {
		identifier   string
		branch       string
		oauthTokenID string
	}

	type inputConfig struct {
		args             []string
		organization     string
		token            string
		workspace        string
		description      string
		workingDirectory string
		vcsConfig        *vcsConfig
	}

	newWorkspaceUpdateOptions := func(
		description string,
		workingDirectory string,
		vcsOpts *vcsConfig,
	) tfe.WorkspaceUpdateOptions {
		r := tfe.WorkspaceUpdateOptions{}
		if description != "" {
			r.Description = &description
		}
		if workingDirectory != "" {
			r.WorkingDirectory = &workingDirectory
		}
		if vcsOpts != nil {
			r.VCSRepo = &tfe.VCSRepoOptions{
				Identifier:   &vcsOpts.identifier,
				Branch:       &vcsOpts.branch,
				OAuthTokenID: &vcsOpts.oauthTokenID,
			}
		}
		return r
	}

	testConfigs := []struct {
		description  string
		inputConfig  inputConfig
		updateConfig *updateConfig
		resultConfig *resultConfig
	}{
		{
			"Updates description",
			inputConfig{
				[]string{"-workspace", "someworkspace", "-description", "new description"},
				"someorg",
				"sometoken",
				"someworkspace",
				"new description",
				"",
				nil,
			},
			&updateConfig{
				resultObj: &tfe.Workspace{
					ID:          "foo",
					Description: "new description",
				},
			},
			&resultConfig{
				resultObj: &CommandResult{
					Result: WorkspacesUpdateCommandResult{
						ID:          "foo",
						Description: "new description",
					},
				},
			},
		},
		{
			"VCS branch not specified",
			inputConfig{
				[]string{"-workspace", "someworkspace", "-vcs-identifier", "someorg/somerepo"},
				"someorg",
				"sometoken",
				"someworkspace",
				"",
				"",
				nil,
			},
			nil,
			&resultConfig{
				resultErr: errors.New("VCS identifier is specified but branch name is not"),
			},
		},
		{
			"VCS OAuth token ID not specified",
			inputConfig{
				[]string{
					"-workspace", "someworkspace", "-vcs-identifier", "someorg/somerepo",
					"-vcs-branch", "somebranch",
				},
				"someorg",
				"sometoken",
				"someworkspace",
				"",
				"",
				nil,
			},
			nil,
			&resultConfig{
				resultErr: errors.New("VCS identifier is specified but OAuth token ID is not"),
			},
		},
		{
			"Updates VCS settings",
			inputConfig{
				[]string{
					"-workspace", "someworkspace", "-vcs-identifier", "someorg/somerepo",
					"-vcs-branch", "somebranch", "-vcs-oauth-token-id", "sometokenid",
				},
				"someorg",
				"sometoken",
				"someworkspace",
				"",
				"",
				&vcsConfig{
					identifier:   "someorg/somerepo",
					branch:       "somebranch",
					oauthTokenID: "sometokenid",
				},
			},
			&updateConfig{
				resultObj: &tfe.Workspace{},
			},
			nil,
		},
		{
			"Updates working directory",
			inputConfig{
				[]string{
					"-workspace", "someworkspace", "-working-directory", "somedirectory",
				},
				"someorg",
				"sometoken",
				"someworkspace",
				"",
				"somedirectory",
				nil,
			},
			&updateConfig{
				resultObj: &tfe.Workspace{},
			},
			nil,
		},
	}
	for _, test := range testConfigs {
		t.Run(test.description, func(t *testing.T) {
			args := append([]string{"workspaces", "update"}, test.inputConfig.args...)
			var stdBuff, errBuff bytes.Buffer
			options := ExecuteOpts{
				AppName: "tfc-cli",
				Stdout:  &stdBuff,
				Stderr:  &errBuff,
			}
			// Set up expectations
			mockedOSProxy := mockOSProxy{}
			mockedOSProxy.On("lookupEnv", "TFC_ORG").Return(test.inputConfig.organization, true)
			mockedOSProxy.On("lookupEnv", "TFC_TOKEN").Return(test.inputConfig.token, true)
			mockedWorkspacesProxy := mockWorkspacesProxy{}

			if test.updateConfig != nil {
				mockedWorkspacesProxy.On(
					"update",
					mock.Anything,
					mock.Anything,
					test.inputConfig.organization,
					test.inputConfig.workspace,
					newWorkspaceUpdateOptions(
						test.inputConfig.description,
						test.inputConfig.workingDirectory,
						test.inputConfig.vcsConfig,
					),
				).Return(
					test.updateConfig.resultObj,
					test.updateConfig.resultErr,
				)
			}

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

			if test.resultConfig != nil {
				if test.resultConfig.resultErr == nil {
					assert.Nil(t, err)
					assert.Empty(t, errBuff.String())
				} else {
					assert.EqualError(t, err, test.resultConfig.resultErr.Error())
					assert.NotEmpty(t, errBuff.String())
				}
				if test.resultConfig.resultObj != nil {
					expectedOutput, _ := json.Marshal(test.resultConfig.resultObj)
					expectedOutput = append(expectedOutput, '\n')
					assert.Equal(t, string(expectedOutput), stdBuff.String())
				} else {
					assert.Empty(t, stdBuff.String())
				}
			}
		})
	}
}
