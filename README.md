# Terraform Cloud CLI

A command line utility for interfacing with Terraform Cloud. Uses [tfe-go][] under the hood.

[![Tests](https://github.com/cbsinteractive/tfc-cli/actions/workflows/tests.yml/badge.svg)](https://github.com/cbsinteractive/tfc-cli/actions/workflows/tests.yml)

## Available Commands

Create a workspace:

```shell
tfc-cli workspaces create -workspace foo
```

Create a workspace variable:

```shell
tfc-cli workspaces variables create -workspace foo -key bar -value baz -category terraform
```

Update a workspace variable value:

```shell
tfc-cli workspaces variables update value -workspace foo -key bar -value quux
```

Make a workspace variable sensitive:

```shell
tfc-cli workspaces variables update sensitive -workspace foo -key bar -sensitive=true
```

Note the `-sensitive=<bool>` syntax. This is required.

Delete a workspace variable:

```shell
tfc-cli workspaces variables delete -workspace foo -key bar
```

Get current state output variable value:

```shell
tfc-cli stateversions current getoutput -workspace foo -name bar
```

[tfe-go]: https://github.com/hashicorp/go-tfe
