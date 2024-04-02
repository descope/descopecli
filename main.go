package main

import (
	"os"

	"github.com/descope/descopecli/accesskey"
	"github.com/descope/descopecli/audit"
	"github.com/descope/descopecli/project"
	"github.com/descope/descopecli/shared"
	"github.com/descope/descopecli/tenant"
	"github.com/descope/descopecli/user"
	"github.com/spf13/cobra"
)

func main() {
	cli := &cobra.Command{
		Version: shared.Version,
		Use:     "descopecli",
		Short:   help,
		Example: examples,
	}

	entity := &cobra.Group{ID: "entity", Title: "Entity Commands:"}
	cli.AddGroup(entity)

	accesskey.AddCommands(cli, entity)
	user.AddCommands(cli, entity)
	tenant.AddCommands(cli, entity)

	proj := &cobra.Group{ID: "proj", Title: "Project Commands:"}
	cli.AddGroup(proj)

	audit.AddCommands(cli, proj)
	project.AddCommands(cli, proj)

	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}

// Documentation

const help = "A command line utility for working with the Descope management APIs"

const examples = `  # Load an existing user by their loginId
  export DESCOPE_PROJECT_ID=...
  export DESCOPE_MANAGEMENT_KEY=...
  descopecli user load -l name@example.com

  # Export all project settings and configurations
  export DESCOPE_MANAGEMENT_KEY=...
  descopecli project snapshot export P2Z1234567890123456789012345`
