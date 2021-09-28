package cmd

import (
	"bytes"
	"testing"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWorkspacesCreate(t *testing.T) {
	newWorkspaceCreateOptions := (func(name string) tfe.WorkspaceCreateOptions {
		description := "Created by tfc-cli"
		return tfe.WorkspaceCreateOptions{
			Name:        &name,
			Description: &description,
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
	}{
		{
			"workspace created",
			[]string{"-workspace", "foo"},
			"some org",
			"some token",
			"foo",
			&tfe.Workspace{},
			nil,
		},
	}
	for _, d := range testConfigs {
		t.Run(d.description, func(t *testing.T) {
			args := append([]string{"workspaces", "create"}, d.args...)
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
			mockedWorkspacesProxy.On("create", mock.Anything, mock.Anything, d.organization, newWorkspaceCreateOptions(d.workspace)).Return(d.workspaceCreateResult, d.workspaceCreateError)

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
			assert.Nil(t, err)
			mockedOSProxy.AssertExpectations(t)
			mockedWorkspacesProxy.AssertExpectations(t)
		})
	}
}
