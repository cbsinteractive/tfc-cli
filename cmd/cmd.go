package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/go-tfe"
)

var Version = "development"

type Usage string

const (
	OrgUsage        Usage = "Organization name"
	TokenUsage      Usage = "Organization token"
	WorkspaceUsage  Usage = "Workspace name"
	OutputNameUsage Usage = "Output variable name"
)

type Runner interface {
	Name() string
	Init([]string) error
	Run() error
}

type OrgOpts struct {
	name  string
	token string
}

type WorkspaceOpts struct {
	name string
}

type ExecuteOpts struct {
	AppName string
	Writer  io.Writer
}

type CommandResult struct {
	Error  string `json:"error,omitempty"`
	Result string `json:"result,omitempty"`
}

func Execute(options ExecuteOpts) error {
	return root(options, os.Args[1:], productionDependencyCaller{}, os.Stdout)
}

func root(options ExecuteOpts, args []string, os dependencyCaller, w io.Writer) error {
	if len(args) < 1 {
		return errors.New("no subcommand given")
	}
	runners := []Runner{
		NewStateVersionsCmd(os, w),
		NewVersionCmd(w),
	}
	subcommand := args[0]
	for _, r := range runners {
		if r.Name() == subcommand {
			if err := r.Init(args[1:]); err != nil {
				return err
			}
			return r.Run()
		}
	}
	return fmt.Errorf("unknown subcommand: %s", subcommand)
}

type dependencyCaller interface {
	osLookupEnv(string) (string, bool)
	clientWorkspacesRead(*tfe.Client, context.Context, string, string) (*tfe.Workspace, error)
	clientStateVersionsCurrentWithOptions(
		*tfe.Client,
		context.Context,
		string,
		*tfe.StateVersionCurrentOptions,
	) (*tfe.StateVersion, error)
}

type productionDependencyCaller struct{}

func (c productionDependencyCaller) osLookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

func (c productionDependencyCaller) clientWorkspacesRead(
	client *tfe.Client,
	ctx context.Context,
	organization string,
	workspace string,
) (*tfe.Workspace, error) {
	return client.Workspaces.Read(ctx, organization, workspace)
}

func (c productionDependencyCaller) clientStateVersionsCurrentWithOptions(
	client *tfe.Client,
	ctx context.Context,
	workspaceID string,
	options *tfe.StateVersionCurrentOptions,
) (*tfe.StateVersion, error) {
	return client.StateVersions.CurrentWithOptions(ctx, workspaceID, options)
}

func processCommonInputs(token *string, orgName *string, os dependencyCaller) error {
	if *token == "" {
		var ok bool
		*token, ok = os.osLookupEnv("TFC_TOKEN")
		if !ok {
			return errors.New(
				"org token must be provided via the -token argument or by setting the TFC_TOKEN environment variable",
			)
		}
	}
	if *orgName == "" {
		var ok bool
		*orgName, ok = os.osLookupEnv("TFC_ORG")
		if !ok {
			return errors.New(
				"org name must be provided via the -org argument or by setting the TFC_ORG environment variable",
			)
		}
	}
	return nil
}
