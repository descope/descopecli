package audit

import (
	"github.com/descope/descopecli/shared"
	"github.com/spf13/cobra"
)

var Flags struct {
	Text string
}

func AddCommands(parent *cobra.Command, group *cobra.Group) {
	audit := shared.MakeGroupCommand(group, "audit", "Commands for working with audit logs")
	parent.AddCommand(audit)

	shared.AddCommand(audit, Search, "search <text>", "Full text search up to last 30 days of audit records", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
	})
}
