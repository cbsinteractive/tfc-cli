package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWorkspacesShow(t *testing.T) {
	testConfigs := []struct {
		description         string
		args                []string
		organization        string
		token               string
		workspace           string
		workspaceShowResult *tfe.Workspace
		workspaceShowError  error
		expectOutput        bool
	}{
		{
			"show existing workspace",
			[]string{"-workspace", "foo"},
			"some org",
			"some token",
			"foo",
			&tfe.Workspace{},
			nil,
			true,
		},
		{
			"show existing workspace (quiet)",
			[]string{"-workspace", "foo", "-quiet"},
			"some org",
			"some token",
			"foo",
			&tfe.Workspace{},
			nil,
			false,
		},
		{
			"show missing workspace",
			[]string{"-workspace", "foo"},
			"some org",
			"some token",
			"foo",
			nil,
			errors.New("resource not found"),
			false,
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
			}
		})
	}
}
