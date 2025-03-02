package apps

import (
	"github.com/descope/descopecli/shared"
	"github.com/spf13/cobra"
)

var Flags struct {
	Description      string
	FlowHostingURL   string
	PermissionScopes []string
	AttributeScopes  []string
	CallbackURLs     []string
}

func AddCommands(parent *cobra.Command, group *cobra.Group) {
	apps := shared.MakeGroupCommand(group, "apps", "Commands for creating and managing applications and integrations")
	parent.AddCommand(apps)

	inbound := shared.MakeGroupCommand(nil, "inbound", "Commands for creating and managing inbound apps")
	apps.AddCommand(inbound)

	shared.AddCommand(inbound, Create, "create <name> <-f flow-hosting-url> <p permission-scopes>", "Create a new inbound app", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
		setCreateFlags(cmd)
	})

	shared.AddCommand(inbound, Update, "update <id> <name> <-f flow-hosting-url> <p permission-scopes>", "Update an existing inbound app", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(2)
		setCreateFlags(cmd)
	})

	shared.AddCommand(inbound, Delete, "delete <id>", "Delete a inbound app", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
	})

	shared.AddCommand(inbound, Load, "load <id>", "Load details about a inbound app", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
	})

	shared.AddCommand(inbound, LoadAll, "load-all", "Load all inbound apps", func(_ *cobra.Command) {
	})

	secret := shared.MakeGroupCommand(nil, "secret", "Commands for managing inbound app secrets")
	inbound.AddCommand(secret)

	shared.AddCommand(secret, LoadSecret, "load <id>", "Load the secret for a inbound app", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
	})

	shared.AddCommand(secret, RotateSecret, "rotate <id>", "Generates a new secret for a inbound app", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
	})
}

func setCreateFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Flags.Description, "description", "", "an optional description for the inbound app")
	cmd.Flags().StringVarP(&Flags.FlowHostingURL, "flow-hosting-url", "f", "", "the flow hosting url for the inbound app")
	cmd.Flags().StringSliceVarP(&Flags.CallbackURLs, "callback-url", "c", nil, "can be used multiple times to add approved callback urls. For example:\n"+"    -c 'https://example.com'")
	cmd.Flags().StringSliceVarP(&Flags.PermissionScopes, "permission-scope", "p", nil, "can be used multiple times to add permission scopes, where each value\n"+"is expected to be a comma separated list with the scope name, description,\n"+"and a list of roles separated by colons. For example:\n"+"    -p 'write,Allow writing files,User|Reader|Writer'\n"+"    -p 'guest,Guest user with no roles,-'")
	cmd.Flags().StringSliceVarP(&Flags.AttributeScopes, "attribute-scope", "a", nil, "can be used multiple times to add attribute scopes, where each value\n"+"is expected to be a comma separated list with the scope name, description,\n"+"and a list of user attributes separated by colons. For example:\n"+"    -a 'contact,Fetch user contact details,displayName|email|phone'")
}
