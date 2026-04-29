# Descope CLI

[![Go](https://img.shields.io/badge/Go-1.21%2B-00ADD8)](#build-from-source)
[![Releases](https://img.shields.io/github/v/release/descope/descopecli)](https://github.com/descope/descopecli/releases/latest)
[![License](https://img.shields.io/github/license/descope/descopecli)](LICENSE)

A command line utility for working with the Descope management APIs. `descope` helps you manage project resources, export and import settings, inspect audit logs, and automate admin workflows from scripts or CI.

## Table of Contents
- [Why use it](#why-use-it)
- [Features](#features)
- [Installation](#installation)
- [Getting started](#getting-started)
- [Common commands](#common-commands)
- [Repository layout](#repository-layout)
- [Build from source](#build-from-source)
- [JSON output for automation](#json-output-for-automation)
- [Contributing](#contributing)

## Why use it

Descope CLI gives teams a repeatable way to manage Descope projects without clicking through the console for every change.

## Features

- Create and modify project resources such as users, tenants, and access keys
- Export, validate, and import project configuration snapshots
- Search and inspect audit logs
- Use JSON output in shell scripts, CI pipelines, and internal tooling
- Work with the same management APIs used by the Descope console

## Installation

### Download a binary

Download the latest prebuilt binary from the [latest release](https://github.com/descope/descopecli/releases/latest).

### Debian or Ubuntu

```bash
sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys e8365d8513142909
echo "deb https://descope.github.io/packages stable main" | sudo tee /etc/apt/sources.list.d/descope.list
sudo apt-get update
sudo apt-get install descope
```

### Fedora or CentOS

```bash
sudo dnf config-manager --add-repo https://descope.github.io/packages/descope.repo
sudo dnf install descope
```

## Getting started

### Requirements

You will need:

- a Descope `Project ID`
- a valid Descope management key for that project

Set both values before running commands:

```bash
export DESCOPE_PROJECT_ID='P...'
export DESCOPE_MANAGEMENT_KEY='K...'
descope --help
```

## Common commands

```bash
descope tenant list
descope user search --help
descope audit search --help
descope project export --help
```

If you prefer explicit arguments, some commands also accept the project id directly.

## Repository layout

- `main.go`, CLI entrypoint
- `tenant/`, tenant-related commands
- `accesskey/`, access key management
- `audit/`, audit log operations
- `project/`, project configuration and snapshots
- `apps/`, `flow/`, `theme/`, additional resource groups
- `shared/`, shared helpers and API clients

## Build from source

Prerequisites:

- Go 1.21+
- Make

```bash
git clone https://github.com/descope/descopecli
cd descopecli
make install
```

This installs the CLI to `$GOPATH/bin` by default.

## JSON output for automation

The CLI supports JSON-friendly workflows, which makes it useful for:

- CI checks
- environment bootstrapping
- project backup and restore
- scripted audits and reporting

Example:

```bash
descope tenant list --help
```

## Contributing

Issues and pull requests are welcome. If you are proposing a larger change, opening an issue first can help align on the expected command UX.

## License

Released under the [MIT License](LICENSE).
