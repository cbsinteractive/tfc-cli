package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
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
	assert.Equal(t, "development", buff.String())
}
