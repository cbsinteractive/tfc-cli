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
			options := ExecuteOpts{
				AppName: "tfc-cli",
			}
			var buff bytes.Buffer
			if err := root(
				options,
				args,
				defaultFakeDeps{},
				&buff); err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, d.expectedOutput, buff.String())
		})
	}
}
