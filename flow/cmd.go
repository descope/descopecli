package flow

import (
	"github.com/descope/descopecli/shared"
	"github.com/spf13/cobra"
)

var Flags struct {
	Skip bool
}

func AddCommands(parent *cobra.Command, group *cobra.Group) {
	flow := shared.MakeGroupCommand(group, "flow", "Commands for managing flows")
	parent.AddCommand(flow)

	shared.AddCommand(flow, RunManagementFlow, "run <flowId>", "Run an autonomous flow", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
	})

	shared.AddCommand(flow, List, "list", "Lists all flows in a project", func(cmd *cobra.Command) {
		cmd.Args = cobra.NoArgs
	})

	shared.AddCommand(flow, Export, "export <flowId> [targetPath]", "Export a flow to a JSON file or standard output", func(cmd *cobra.Command) {
		cmd.Args = cobra.RangeArgs(1, 2)
	})

	shared.AddCommand(flow, Import, "import <flowId> <sourcePath>", "Import a flow from a JSON file", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(2)
	})

	shared.AddCommand(flow, Convert, "convert <sourcePath> [targetPath] [-s]", "Convert a flow between formats", func(cmd *cobra.Command) {
		cmd.Args = cobra.RangeArgs(1, 2)
		cmd.Flags().BoolVarP(&Flags.Skip, "skip", "s", false, "skip unsupported file types")
		cmd.PreRunE = shared.StandalonePreRun
	})
}
