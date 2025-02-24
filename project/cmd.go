package project

import (
	"github.com/descope/descopecli/shared"
	"github.com/descope/descopecli/snapshot"
	"github.com/spf13/cobra"
)

var Flags struct {
	Environment string
	Tags        []string
	Force       bool
}

func AddCommands(parent *cobra.Command, group *cobra.Group) {
	project := shared.MakeGroupCommand(group, "project", "Commands for managing projects")
	parent.AddCommand(project)

	shared.AddCommand(project, List, "list", "Lists all projects in a company", func(cmd *cobra.Command) {
		cmd.Args = cobra.NoArgs
	})

	shared.AddCommand(project, Clone, "clone <existingProjectId> <newProjectName> [-e environment] [--tags tag,...]", "Clone an existing project along with all settings and configurations", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(2)
		cmd.Flags().StringVarP(&Flags.Environment, "environment", "e", "", "an optional environment for the new project, only valid value is production")
		cmd.Flags().StringSliceVar(&Flags.Tags, "tags", nil, "a comma separated list of tags for the new project")
		cmd.PreRunE = shared.ProjectPreRun

		// deprecated
		cmd.Flags().StringVarP(&Flags.Environment, "tag", "t", "", "an optional tag for the new project, only valid value is production")
		cmd.Flags().MarkDeprecated("tag", "use --environment instead")
		cmd.Flags().MarkShorthandDeprecated("t", "use -e instead")
	})

	shared.AddCommand(project, Delete, "delete <projectId> [-f]", "Delete an existing project", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
		cmd.Flags().BoolVarP(&Flags.Force, "force", "f", false, "skips the prompt and deletes the project immediately")
		cmd.PreRunE = shared.ProjectPreRun
	})

	snap := shared.MakeGroupCommand(nil, "snapshot", "Commands for working with project snapshots")
	// next version: snap.Deprecated = `Use "descope snapshot" instead of "descope project snapshot".`
	project.AddCommand(snap)
	snapshot.AddSnapshotCommands(snap)
}
