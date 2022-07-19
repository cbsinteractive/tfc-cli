# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- New `-terraformVersion` option to `workspaces create` subcommand.

## [1.1.0] - 2022-04-25

### Added

- New `workspaces set-auto-apply` subcommand.

### Changed

- Update dependency `go-tfe` to v1.1.0.

## [1.0.0] - 2022-01-24

### Added

- New `workspaces set-description` subcommand.
- New `workspaces set-working-directory` subcommand.
- New `workspaces set-vcs-branch` subcommand.
- New `workspaces unset-description` subcommand.
- New `workspaces unset-working-directory` subcommand.
- New `workspaces unset-vcs-branch` subcommand.

### Removed

- Removed `workspaces update` subcommand in favor of individual set/unset subcommands.
- Removed meta settings from `workspaces create` command in favor of individual set/unset subcommands. The user should call the various set/unset subcommands after the workspace has been created.
