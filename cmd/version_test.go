package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	testConfigs := []struct {
		description    string
		major          string
		minor          string
		patch          string
		label          string
		expectedOutput string
	}{
		{"default", "1", "2", "3", "", "1.2.3"},
		{"includes label", "1", "2", "3", "foo", "1.2.3-foo"},
	}
	for _, d := range testConfigs {
		t.Run(d.description, func(t *testing.T) {
			Major = d.major
			Minor = d.minor
			Patch = d.patch
			ReleaseLabel = d.label
			args := []string{"version"}
			var buff bytes.Buffer
			options := ExecuteOpts{
				AppName: "tfc-cli",
				Writer:  &buff,
			}

			// Code under test
			err := root(
				options,
				args,
				dependencyProxies{},
			)

			// Verify
			assert.Nil(t, err)
			assert.Equal(t, d.expectedOutput, buff.String())
		})
	}
}
