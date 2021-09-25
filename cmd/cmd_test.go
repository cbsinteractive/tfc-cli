package cmd

import (
	"context"
	"errors"

	"github.com/hashicorp/go-tfe"
)

type defaultFakeDeps struct{}

func (c defaultFakeDeps) osLookupEnv(key string) (string, bool) {
	return "", false
}

func (c defaultFakeDeps) clientWorkspacesRead(
	_ *tfe.Client,
	ctx context.Context,
	organization string,
	workspace string,
) (*tfe.Workspace, error) {
	return nil, errors.New("not implemented")
}

func (c defaultFakeDeps) clientStateVersionsCurrentWithOptions(
	_ *tfe.Client,
	ctx context.Context,
	workspaceID string,
	options *tfe.StateVersionCurrentOptions,
) (*tfe.StateVersion, error) {
	return nil, errors.New("not implemented")
}
