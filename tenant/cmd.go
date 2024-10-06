package tenant

import (
	"github.com/descope/descopecli/shared"
	"github.com/spf13/cobra"
)

var Flags struct {
	TenantID string
	Name     string
	Domains  []string
}

func AddCommands(parent *cobra.Command, group *cobra.Group) {
	tenant := shared.MakeGroupCommand(group, "tenant", "Commands for creating and managing tenants")
	parent.AddCommand(tenant)

	shared.AddCommand(tenant, Create, "create <name>", "Create a new tenant", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
		cmd.Flags().StringVarP(&Flags.TenantID, "id", "i", "", "an optional custom id for the tenant")
		cmd.Flags().StringSliceVarP(&Flags.Domains, "domains", "d", nil, "a comma separated list of self provisioning domains for the tenant")
	})

	shared.AddCommand(tenant, Delete, "delete <id>", "Delete a tenant", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
	})

	shared.AddCommand(tenant, Load, "load <id>", "Load details about a tenant", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
	})

	shared.AddCommand(tenant, LoadAll, "load-all", "Load all tenants", func(_ *cobra.Command) {
	})
}
