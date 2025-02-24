package snapshot

import (
	"github.com/descope/descopecli/shared"
	"github.com/spf13/cobra"
)

var Flags struct {
	Path           string
	Debug          bool
	SecretsInput   string
	SecretsOutput  string
	FailuresOutput string
	NoAssets       bool
}

func AddCommands(parent *cobra.Command, group *cobra.Group) {
	snapshot := shared.MakeGroupCommand(group, "snapshot", "Commands for working with snapshots")
	parent.AddCommand(snapshot)

	AddSnapshotCommands(snapshot)
}

func AddSnapshotCommands(snapshot *cobra.Command) {
	shared.AddCommand(snapshot, Export, "export <projectId> [-p path]", "Export a snapshot of all the settings of a project", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
		cmd.Flags().StringVarP(&Flags.Path, "path", "p", "", "the path to write the snapshot into")
		cmd.Flags().BoolVar(&Flags.NoAssets, "no-assets", false, "don't extract assets from snapshot files")
		cmd.PreRunE = shared.ProjectPreRun
	})

	shared.AddCommand(snapshot, Import, "import <projectId> [-p path]", "Import a snapshot into a project", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
		cmd.Flags().StringVarP(&Flags.Path, "path", "p", "", "the path to read the snapshot from")
		cmd.Flags().StringVar(&Flags.SecretsInput, "secrets-input", "", "the path to a JSON file with required secrets")
		cmd.PreRunE = shared.ProjectPreRun
	})

	shared.AddCommand(snapshot, Validate, "validate <projectId> [-p path]", "Validate a snapshot before importing into a project", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
		cmd.Flags().StringVarP(&Flags.Path, "path", "p", "", "the path to read the snapshot from")
		cmd.Flags().StringVar(&Flags.SecretsInput, "secrets-input", "", "the path to a JSON file with required secrets")
		cmd.Flags().StringVar(&Flags.SecretsOutput, "secrets-output", "", "the path to a JSON file to write missing secrets in case validation fails")
		cmd.Flags().StringVar(&Flags.FailuresOutput, "failures-output", "", "the path to write a list of failures in case validation fails")
		cmd.PreRunE = shared.ProjectPreRun
	})
}
