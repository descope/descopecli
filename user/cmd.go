package user

import (
	"github.com/descope/descopecli/shared"
	"github.com/spf13/cobra"
)

var Flags struct {
	LoginID string
	UserID  string
	Email   string
	Phone   string
	Name    string
	Tenants []string
	Limit   int
	Page    int
}

func AddCommands(parent *cobra.Command, group *cobra.Group) {
	user := shared.MakeGroupCommand(group, "user", "Commands for creating and managing users")

	shared.AddCommand(user, Create, "create <loginId> [-e email] [-p phone] [-n name] [-t tid,...]", "Create a new user", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
		cmd.Flags().StringVarP(&Flags.Email, "email", "e", "", "the user's email address")
		cmd.Flags().StringVarP(&Flags.Phone, "phone", "p", "", "the user's phone number")
		cmd.Flags().StringVarP(&Flags.Name, "name", "n", "", "the user's display name")
		cmd.Flags().StringSliceVarP(&Flags.Tenants, "tenants", "t", nil, "a comma separated list of tenant ids for the user")
	})

	shared.AddCommand(user, Delete, "delete {-l loginId | -u userId}", "Delete a user", func(cmd *cobra.Command) {
		cmd.Flags().StringVarP(&Flags.LoginID, "login-id", "l", "", "the user's loginId")
		cmd.Flags().StringVarP(&Flags.UserID, "user-id", "u", "", "the user's userId")
	})

	shared.AddCommand(user, Load, "load {-l loginId | -u userId}", "Load details about a user", func(cmd *cobra.Command) {
		cmd.Flags().StringVarP(&Flags.LoginID, "login-id", "l", "", "the user's loginId")
		cmd.Flags().StringVarP(&Flags.UserID, "user-id", "u", "", "the user's userId")
	})

	shared.AddCommand(user, LoadAll, "load-all [-l limit] [-p page]", "Load all users", func(cmd *cobra.Command) {
		cmd.Flags().IntVarP(&Flags.Limit, "limit", "l", 100, "the number of results for pagination (max 100)")
		cmd.Flags().IntVarP(&Flags.Page, "page", "p", 0, "the number of page for pagination (default 0)")
	})

	shared.AddCommand(user, SetTemporaryPassword, "set-temporary-password <loginId> <password>", "Set a temporary password for a user", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(2)
	})

	shared.AddCommand(user, SetActivePassword, "set-active-password <loginId> <password>", "Set an active password for a user", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(2)
	})

	shared.AddCommand(user, ExpirePassword, "expire-password <loginId>", "Expire a user's password", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
	})

	parent.AddCommand(user)
}
