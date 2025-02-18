package thirdpartyapps

import (
	"github.com/descope/descopecli/shared"
	"github.com/spf13/cobra"
)

var Flags struct {
	Description      string
	PermissionScopes []string
	AttributeScopes  []string
	CallbackURLs     []string
}

func AddCommands(parent *cobra.Command, group *cobra.Group) {
	thirdpartyapps := shared.MakeGroupCommand(group, "thirdpartyapps", "Commands for creating and managing third party applications")
	parent.AddCommand(thirdpartyapps)

	shared.AddCommand(thirdpartyapps, Create, "create <name> <flow-hosting-url>", "Create a new third party app", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(2)
		setCreateFlags(cmd)
	})

	shared.AddCommand(thirdpartyapps, Update, "update <id> <name> <flow-hosting-url>", "Update the details for a third party app", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(3)
		setCreateFlags(cmd)
	})

	shared.AddCommand(thirdpartyapps, Delete, "delete <id>", "Delete a third party app", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
	})

	shared.AddCommand(thirdpartyapps, Load, "load <id>", "Load details about a third party app", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
	})

	shared.AddCommand(thirdpartyapps, LoadAll, "load-all", "Load all third party apps", func(_ *cobra.Command) {
	})

	secret := shared.MakeGroupCommand(nil, "secret", "Commands for managing secrets")
	thirdpartyapps.AddCommand(secret)

	shared.AddCommand(secret, LoadSecret, "load <id>", "Load the secret for a third party app", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
	})

	shared.AddCommand(secret, RotateSecret, "rotate <id>", "Generates a new secret for a third party app", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
	})
}

func setCreateFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Flags.Description, "description", "", "an optional description for the third party app")
	cmd.Flags().StringSliceVarP(&Flags.CallbackURLs, "callback-url", "c", nil, "can be used multiple times to add approved callback urls. For example:\n"+"    -c 'https://example.com'")
	cmd.Flags().StringSliceVarP(&Flags.PermissionScopes, "permission-scopes", "p", nil, "can be used multiple times to add permission scopes, where each value\n"+"is expected to be a comma separated list with the scope name, description,\n"+"and a list of roles separated by colons. For example:\n"+"    -p 'write,Allow writing files,User|Reader|Writer'\n"+"    -p 'guest,Guest user with no roles,-'")
	cmd.Flags().StringSliceVarP(&Flags.AttributeScopes, "attribute-scopes", "a", nil, "can be used multiple times to add attribute scopes, where each value\n"+"is expected to be a comma separated list with the scope name, description,\n"+"and a list of user attributes separated by colons. For example:\n"+"    -a 'contact,Fetch user contact details,displayName|email|phone'")
}
