package cmd

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
)

type testStateVersionsCurrentGetOutputDeps struct {
	values      map[string]string
	workspaceId string
	outputs     []*tfe.StateVersionOutput
}

func (c testStateVersionsCurrentGetOutputDeps) osLookupEnv(key string) (string, bool) {
	if _, ok := c.values[key]; !ok {
		return "", false
	}
	return c.values[key], true
}

func (c testStateVersionsCurrentGetOutputDeps) clientWorkspacesRead(
	_ *tfe.Client,
	ctx context.Context,
	organization string,
	workspace string,
) (*tfe.Workspace, error) {
	if c.workspaceId == "" {
		return nil, errors.New("resource not found")
	}
	return &tfe.Workspace{
		ID: c.workspaceId,
	}, nil
}

func (c testStateVersionsCurrentGetOutputDeps) clientStateVersionsCurrentWithOptions(
	_ *tfe.Client,
	ctx context.Context,
	workspaceID string,
	options *tfe.StateVersionCurrentOptions,
) (*tfe.StateVersion, error) {
	if c.outputs == nil {
		return nil, errors.New("not implemented")
	}
	return &tfe.StateVersion{
		Outputs: c.outputs,
	}, nil
}

func TestStateVersionsCurrentGetOutput(t *testing.T) {
	defaultEnv := func() map[string]string {
		return map[string]string{
			"TFC_TOKEN": "some token",
			"TFC_ORG":   "some org",
		}
	}
	testConfigs := []struct {
		name          string
		env           map[string]string
		workspaceId   string
		outputName    string
		outputs       []*tfe.StateVersionOutput
		expectedValue string
	}{
		{
			"output variable found",
			defaultEnv(),
			"some workspace id",
			"foo",
			[]*tfe.StateVersionOutput{
				{
					Name:  "foo",
					Value: "some value",
				},
			},
			"some value"},
	}
	for _, d := range testConfigs {
		t.Run(d.name, func(t *testing.T) {
			args := []string{
				"stateversions",
				"current",
				"getoutput",
				"-workspace",
				"some workspace",
				"-name",
				d.outputName,
			}
			options := ExecuteOpts{
				AppName: "tfc-cli",
			}
			var buff bytes.Buffer
			if err := root(
				options,
				args,
				testStateVersionsCurrentGetOutputDeps{
					values:      d.env,
					workspaceId: d.workspaceId,
					outputs:     d.outputs,
				},
				&buff); err != nil {
				t.Fatal(err)
			}
			assert.Contains(t, buff.String(), d.expectedValue)
		})
	}
}
