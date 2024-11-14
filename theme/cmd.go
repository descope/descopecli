package theme

import (
	"github.com/descope/descopecli/shared"
	"github.com/spf13/cobra"
)

var Flags struct {
	Skip bool
}

func AddCommands(parent *cobra.Command, group *cobra.Group) {
	theme := shared.MakeGroupCommand(group, "theme", "Commands for managing themes")
	parent.AddCommand(theme)

	shared.AddCommand(theme, Export, "export [targetPath]", "Export the project theme to a JSON file or standard output", func(cmd *cobra.Command) {
		cmd.Args = cobra.RangeArgs(0, 1)
	})

	shared.AddCommand(theme, Import, "import <sourcePath>", "Import the project theme from a JSON file", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
	})
}
