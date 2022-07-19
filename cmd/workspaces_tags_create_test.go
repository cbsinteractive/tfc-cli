package cmd

import (
	"bytes"
	"testing"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWorkspacesTagsCreate(t *testing.T) {
	testConfigs := []struct {
		description   string
		args          []string
		organization  string
		token         string
		workspace     string
		workspaceID   string
		tag           string
		expectedError error
	}{
		{
			"add a tag",
			[]string{
				"-workspace",
				"foo",
				"-tag",
				"bar",
			},
			"some org",
			"some token",
			"foo",
			"some workspace id",
			"bar",
			nil,
		},
	}
	for _, d := range testConfigs {
		t.Run(d.description, func(t *testing.T) {
			args := append([]string{"workspaces", "tags", "create"}, d.args...)
			var buff bytes.Buffer
			options := ExecuteOpts{
				AppName: "tfc-cli",
				Stdout:  &buff,
			}
			// Set up expectations
			mockedOSProxy := mockOSProxy{}
			mockedOSProxy.On("lookupEnv", "TFC_ORG").Return(d.organization, true)
			mockedOSProxy.On("lookupEnv", "TFC_TOKEN").Return(d.token, true)
			mockWorkspacesProxy := mockWorkspacesProxy{}
			mockWorkspacesProxy.On("read", mock.Anything, mock.Anything, d.organization, d.workspace).Return(&tfe.Workspace{ID: d.workspaceID}, nil)
			mockedTagsProxy := mockWorkspacesTagsProxy{}
			mockedTagsProxy.On(
				"create",
				mock.Anything,
				mock.Anything,
				d.workspaceID,
				tfe.WorkspaceAddTagsOptions{
					Tags: []*tfe.Tag{
						{
							Name: d.tag,
						},
					},
				}).Return(d.expectedError)

			// Code under test
			err := root(
				options,
				args,
				dependencyProxies{
					client: clientProxy{
						workspaces: mockWorkspacesProxy,
						workspacesCommands: workspacesCommands{
							tags: mockedTagsProxy,
						},
					},
					os: mockedOSProxy,
				},
			)

			// Verify
			assert.Nil(t, err)
			mockedOSProxy.AssertExpectations(t)
			mockWorkspacesProxy.AssertExpectations(t)
			mockedTagsProxy.AssertExpectations(t)
		})
	}
}
