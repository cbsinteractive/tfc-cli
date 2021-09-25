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
		createdWorkspace *tfe.Workspace
		createError      error
	}{
		{"foo", []string{"-workspace", "foo"}, &tfe.Workspace{}, nil},
	}
	for _, d := range testConfigs {
		t.Run(d.description, func(t *testing.T) {
			args := append([]string{"workspaces", "delete"}, d.args...)
			var buff bytes.Buffer
			options := ExecuteOpts{
				AppName: "tfc-cli",
				Writer:  &buff,
			}
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
					os: osProxyForTests{
						envVars: newDefaultEnvForTests(),
					},
				},
			); err != nil {
				t.Fatal(err)
			}
		})
	}
}
