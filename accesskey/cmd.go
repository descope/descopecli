package accesskey

import (
	"github.com/descope/descopecli/shared"
	"github.com/spf13/cobra"
)

var Flags struct {
	Description string
	Expires     int64
	UserId      string
	Tenants     []string
}

func AddCommands(parent *cobra.Command, group *cobra.Group) {
	accessKey := shared.MakeGroupCommand(group, "access-key", "Commands for creating and managing access keys")
	parent.AddCommand(accessKey)

	shared.AddCommand(accessKey, Create, "create <name> [-d description] [-e time] [-t tenants] [-u userId]", "Create a new access key", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
		cmd.Flags().StringVarP(&Flags.Description, "description", "d", "", "an optional description for the access key")
		cmd.Flags().Int64VarP(&Flags.Expires, "expires", "e", 0, "the access key's expiry time (unix time in seconds, default 0 to not expire)")
		cmd.Flags().StringSliceVarP(&Flags.Tenants, "tenants", "t", nil, "a comma separated list of tenant ids for the access key")
		cmd.Flags().StringVarP(&Flags.UserId, "userId", "u", "", "an optional user id to adopt authorizations from")
	})

	shared.AddCommand(accessKey, Delete, "delete <id>", "Delete an access key", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
	})

	shared.AddCommand(accessKey, Load, "load <id>", "Load details about an access key", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
	})

	shared.AddCommand(accessKey, LoadAll, "load-all", "Load all access keys", func(_ *cobra.Command) {
	})

	shared.AddCommand(accessKey, Activate, "activate <id>", "Activate an access key", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
	})

	shared.AddCommand(accessKey, Deactivate, "deactivate <id>", "Deactivate an access key", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
	})
}
