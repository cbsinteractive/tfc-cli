package cmd

import (
	"bytes"
	"testing"

	"github.com/hashicorp/go-tfe"
)

func TestWorkspacesDelete(t *testing.T) {
	testConfigs := []struct {
		description      string
		args             []string
		organization     string
		token            string
		createdWorkspace *tfe.Workspace
		createError      error
	}{
		{
			"foo",
			[]string{"-workspace", "foo"},
			"some org",
			"some token",
			&tfe.Workspace{},
			nil,
		},
	}
	for _, d := range testConfigs {
		t.Run(d.description, func(t *testing.T) {
			args := append([]string{"workspaces", "delete"}, d.args...)
			var buff bytes.Buffer
			options := ExecuteOpts{
				AppName: "tfc-cli",
				Writer:  &buff,
			}
			// Set up expectations
			mockedOSProxy := mockOSProxy{}
			mockedOSProxy.On("lookupEnv", "TFC_ORG").Return(d.organization, true)
			mockedOSProxy.On("lookupEnv", "TFC_TOKEN").Return(d.token, true)
			if err := root(
				options,
				args,
				dependencyProxies{
					client: clientProxy{
						workspaces: workspacesProxyForTests{
							t:            t,
							organization: "some org",
							workspace:    "foo",
						},
					},
					os: mockedOSProxy,
				},
			); err != nil {
				t.Fatal(err)
			}
		})
	}
}
