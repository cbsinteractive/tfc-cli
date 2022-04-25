package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/go-tfe"
)

func variableFromKey(client *tfe.Client, proxy clientProxy, ctx context.Context, workspaceID string, key string) (*tfe.Variable, error) {
	v, err := proxy.workspacesCommands.variables.list(client, ctx, workspaceID, &tfe.VariableListOptions{})
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

type CommandResult struct {
	Result interface{} `json:"result,omitempty"`
}

func output(w io.Writer, result interface{}) {
	r := CommandResult{
		Result: result,
	}
	d, _ := json.Marshal(r)
	d = append(d, '\n')
	w.Write(d)
}

type CommandError struct {
	Error string `json:"error,omitempty"`
}

func outputError(w io.Writer, err error) {
	r := CommandError{
		Error: err.Error(),
	}
	d, _ := json.Marshal(r)
	d = append(d, '\n')
	w.Write(d)
}
