package user

import (
	"github.com/descope/descopecli/shared"
	"github.com/spf13/cobra"
)

var Flags struct {
	LoginID     string
	UserID      string
	TenantID    string
	Email       string
	Phone       string
	Name        string
	Tenants     []string
	Roles       []string
	Limit       int
	Page        int
	RedirectURL string
}

// reused in create commands for regular users and test users
var createUse = "create <loginId> [-e email] [-p phone] [-n name] [-t tid,...]"
var createSetup = func(cmd *cobra.Command) {
	cmd.Args = cobra.ExactArgs(1)
	cmd.Flags().StringVarP(&Flags.Email, "email", "e", "", "the user's email address")
	cmd.Flags().StringVarP(&Flags.Phone, "phone", "p", "", "the user's phone number")
	cmd.Flags().StringVarP(&Flags.Name, "name", "n", "", "the user's display name")
	cmd.Flags().StringSliceVarP(&Flags.Tenants, "tenants", "t", nil, "a comma separated list of tenant ids for the user")
}

func AddCommands(parent *cobra.Command, group *cobra.Group) {
	user := shared.MakeGroupCommand(group, "user", "Commands for creating and managing users")
	parent.AddCommand(user)

	// reused in create user and create test user
	shared.AddCommand(user, Create, createUse, "Create a new user", createSetup)

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

	shared.AddCommand(user, Activate, "activate <loginId>", "Activate a user", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
	})

	shared.AddCommand(user, Deactivate, "deactivate <loginId>", "Deactivate a user", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
	})

	pwd := shared.MakeGroupCommand(nil, "password", "Commands for managing user passwords")
	user.AddCommand(pwd)

	shared.AddCommand(pwd, SetTemporaryPassword, "set-temporary <loginId> <password>", "Set a temporary password for a user", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(2)
	})

	shared.AddCommand(pwd, SetActivePassword, "set-active <loginId> <password>", "Set an active password for a user", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(2)
	})

	shared.AddCommand(pwd, ExpirePassword, "expire <loginId>", "Expire a user's password", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
	})

	roles := shared.MakeGroupCommand(nil, "roles", "Commands for managing user roles")
	user.AddCommand(roles)

	shared.AddCommand(roles, SetRoles, "set <loginId> [-t tid] [-r role,...]", "Set the roles for a user", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
		cmd.Flags().StringSliceVarP(&Flags.Roles, "roles", "r", nil, "a comma separated list of role names, or remove all roles if not set")
		cmd.Flags().StringVarP(&Flags.TenantID, "tenant", "t", "", "update the roles for the user in a specific tenant")
	})

	shared.AddCommand(roles, AddRoles, "add <loginId> [-t tid] <-r role,...>", "Add roles to a user", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
		cmd.Flags().StringSliceVarP(&Flags.Roles, "roles", "r", nil, "a comma separated list of role names to add")
		cmd.Flags().StringVarP(&Flags.TenantID, "tenant", "t", "", "update the roles for the user in a specific tenant")
	})

	shared.AddCommand(roles, RemoveRoles, "remove <loginId> [-t tid] <-r role,...>", "Remove roles from a user", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
		cmd.Flags().StringSliceVarP(&Flags.Roles, "roles", "r", nil, "a comma separated list of role names to remove")
		cmd.Flags().StringVarP(&Flags.TenantID, "tenant", "t", "", "update the roles for the user in a specific tenant")
	})

	test := shared.MakeGroupCommand(nil, "test", "Commands for creating and managing test users")
	user.AddCommand(test)

	shared.AddCommand(test, CreateTestUser, createUse, "Create a new test user", createSetup)

	shared.AddCommand(test, DeleteAllTestUsers, "delete-all", "Delete all existing test users in the project", func(cmd *cobra.Command) {
	})
}
