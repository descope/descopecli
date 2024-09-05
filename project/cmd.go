package project

import (
	"github.com/descope/descopecli/shared"
	"github.com/spf13/cobra"
)

var Flags struct {
	Path           string
	Debug          bool
	Environment    string
	Tags           []string
	SecretsInput   string
	SecretsOutput  string
	FailuresOutput string
	Force          bool
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

	snapshot := shared.MakeGroupCommand(nil, "snapshot", "Commands for working with project snapshots")
	project.AddCommand(snapshot)

	shared.AddCommand(snapshot, Export, "export <projectId> [-p path]", "Export a snapshot of all the settings and configurations of a project", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
		cmd.Flags().StringVarP(&Flags.Path, "path", "p", "", "the path to write the snapshot into")
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
