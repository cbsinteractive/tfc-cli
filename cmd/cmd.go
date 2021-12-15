package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

// Variables set at build time used to generate the version number
var (
	Major        string = "0"
	Minor        string = "0"
	Patch        string = "0"
	ReleaseLabel string
)

type Usage string

const (
	OrgUsage                 Usage = "Organization name"
	TokenUsage               Usage = "Organization token"
	WorkspaceUsage           Usage = "Workspace name"
	OutputNameUsage          Usage = "Output variable name"
	QuietUsage               Usage = "Quiet output"
	VariableKeyUsage         Usage = "Variable key"
	VariableValueUsage       Usage = "Variable value"
	VariableDescriptionUsage Usage = "Variable description"
	VariableCategoryUsage    Usage = "Variable category"
	VariableHCLUsage         Usage = "Indicates whether variable value is HCL-formatted"
	VariableSensitiveUsage   Usage = "Indicates whether variable is sensitive"
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
	Stdout  io.Writer
	Stderr  io.Writer
}

type commonOpts struct {
	quiet bool
}

type CommandResult struct {
	Error  string      `json:"error,omitempty"`
	Result interface{} `json:"result,omitempty"`
}

func Execute(options ExecuteOpts) error {
	return root(
		options,
		os.Args[1:],
		dependencyProxies{
			client: clientProxy{
				stateVersions:      newStateVersionsProxy(),
				workspaces:         newWorkspacesProxy(),
				workspacesCommands: newWorkspacesCommands(),
			},
			os: newOSProxy(),
		},
	)
}

func root(options ExecuteOpts, args []string, deps dependencyProxies) error {
	if len(args) < 1 {
		err := errors.New("no subcommand given")
		outputError(options.Stderr, err)
		return err
	}
	runners := []Runner{
		NewStateVersionsCmd(deps, options.Stdout),
		NewVersionCmd(options.Stdout),
		NewWorkspacesCmd(options, deps),
	}
	subcommand := args[0]
	for _, r := range runners {
		if r.Name() == subcommand {
			if err := r.Init(args[1:]); err != nil {
				outputError(options.Stderr, err)
				return err
			}
			if err := r.Run(); err != nil {
				outputError(options.Stderr, err)
				return err
			}
			return nil
		}
	}
	err := fmt.Errorf("unknown subcommand: %s", subcommand)
	outputError(options.Stderr, err)
	return err
}

func setCommonFlagsetOptions(fs *flag.FlagSet, o *OrgOpts, w *WorkspaceOpts) {
	fs.StringVar(&o.name, "org", "", string(OrgUsage))
	fs.StringVar(&o.token, "token", "", string(TokenUsage))
	fs.StringVar(&w.name, "workspace", "", string(WorkspaceUsage))
}

func processSubcommand(childRunner *Runner, args []string, childRunners []Runner) error {
	subcommand := args[0]
	for _, r := range childRunners {
		if r.Name() == subcommand {
			if err := r.Init(args[1:]); err != nil {
				return err
			}
			*childRunner = r
			return nil
		}
	}
	return fmt.Errorf("unexpected subcommand: %s", subcommand)
}

func processCommonInputs(token *string, orgName *string, lookupEnv func(string) (string, bool)) error {
	if *token == "" {
		var ok bool
		*token, ok = lookupEnv("TFC_TOKEN")
		if !ok {
			return errors.New(
				"org token must be provided via the -token argument or by setting the TFC_TOKEN environment variable",
			)
		}
	}
	if *orgName == "" {
		var ok bool
		*orgName, ok = lookupEnv("TFC_ORG")
		if !ok {
			return errors.New(
				"org name must be provided via the -org argument or by setting the TFC_ORG environment variable",
			)
		}
	}
	return nil
}
