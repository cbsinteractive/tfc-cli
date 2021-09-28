package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-tfe"
)

func variableFromKey(client *tfe.Client, proxy clientProxy, ctx context.Context, workspaceID string, key string) (*tfe.Variable, error) {
	v, err := proxy.workspacesCommands.variables.list(client, ctx, workspaceID, tfe.VariableListOptions{})
	if err != nil {
		return nil, err
	}
	for _, i := range v.Items {
		if i.Key == key {
			return i, nil
		}
	}
	return nil, fmt.Errorf("variable %s not found", key)
}

func newCommandResultOutput(v interface{}) []byte {
	d, _ := json.Marshal(CommandResult{
		Result: v,
	})
	return append(d, '\n')
}

func newCommandErrorOutput(err error) []byte {
	d, _ := json.Marshal(CommandResult{
		Error: err.Error(),
	})
	return append(d, '\n')
}
