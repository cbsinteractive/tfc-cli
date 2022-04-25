package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWorkspacesSetAutoApply(t *testing.T) {
	type inputConfig struct {
		args         []string
		organization string
		token        string
		workspace    string
		autoApply    bool
	}

	type updateConfig struct {
		resultObj *tfe.Workspace
		resultErr error
	}

	type resultConfig struct {
		resultObj *CommandResult
		resultErr error
	}

	newWorkspaceUpdateOptions := func(
		autoApply bool,
	) tfe.WorkspaceUpdateOptions {
		return tfe.WorkspaceUpdateOptions{
			AutoApply: &autoApply,
		}
	}

	testConfigs := []struct {
		description  string
		inputConfig  inputConfig
		updateConfig *updateConfig
		resultConfig *resultConfig
	}{
		{
			"Updates AutoApply to true",
			inputConfig{
				[]string{"-workspace", "someworkspace", "-auto-apply=true"},
				"someorg",
				"sometoken",
				"someworkspace",
				true,
			},
			&updateConfig{
				resultObj: &tfe.Workspace{
					ID:          "foo",
					Description: "some description",
				},
			},
			&resultConfig{
				resultObj: &CommandResult{
					Result: WorkspacesUpdateCommandResult{
						ID:          "foo",
						Description: "some description",
					},
				},
			},
		},
		{
			"Updates AutoApply to true",
			inputConfig{
				[]string{"-workspace", "someworkspace", "-auto-apply=false"},
				"someorg",
				"sometoken",
				"someworkspace",
				false,
			},
			&updateConfig{
				resultObj: &tfe.Workspace{
					ID:          "foo",
					Description: "some description",
				},
			},
			&resultConfig{
				resultObj: &CommandResult{
					Result: WorkspacesUpdateCommandResult{
						ID:          "foo",
						Description: "some description",
					},
				},
			},
		},
	}

	for _, test := range testConfigs {
		t.Run(test.description, func(t *testing.T) {
			args := append([]string{"workspaces", "set-auto-apply"}, test.inputConfig.args...)
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
						test.inputConfig.autoApply,
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
