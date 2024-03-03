package project

import (
	"github.com/descope/descopecli/shared"
	"github.com/spf13/cobra"
)

var Flags struct {
	Path    string
	Debug   bool
	Tag     string
	Secrets string
	Output  string
	Format  string
}

func AddCommands(parent *cobra.Command, group *cobra.Group) {
	project := shared.MakeGroupCommand(group, "project", "Commands for managing project environments")

	shared.AddCommand(project, Export, "export <projectId> [-p path]", "Export all project settings and configurations", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
		cmd.Flags().StringVarP(&Flags.Path, "path", "p", "", "The path to write the project data into")
		cmd.Flags().BoolVar(&Flags.Debug, "debug", false, "Saves an export.log trace file in the debug directory")
		cmd.PreRunE = shared.ProjectPreRun
	})

	shared.AddCommand(project, Import, "import <projectId> [-p path] [-s secrets]", "Import settings and configurations and apply them to a project", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
		cmd.Flags().StringVarP(&Flags.Path, "path", "p", "", "The path to read the project data from")
		cmd.Flags().StringVarP(&Flags.Secrets, "secrets", "s", "", "The path to a JSON file with required secrets")
		cmd.Flags().BoolVar(&Flags.Debug, "debug", false, "Saves an import.log trace file in the debug directory")
		cmd.PreRunE = shared.ProjectPreRun
	})

	shared.AddCommand(project, Validate, "validate <projectId> [-p path] [-s secrets]", "Validate project settings and configurations before import", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
		cmd.Flags().StringVarP(&Flags.Path, "path", "p", "", "The path to read the project data from")
		cmd.Flags().StringVarP(&Flags.Secrets, "secrets", "s", "", "The path to a JSON file with required secrets")
		cmd.Flags().StringVarP(&Flags.Output, "output", "o", "", "The path to write the validation details in case it fails")
		cmd.Flags().StringVarP(&Flags.Format, "format", "f", "", "The validation details output format (json or action, default json)")
		cmd.PreRunE = shared.ProjectPreRun
	})

	shared.AddCommand(project, Clone, "clone <existingProjectId> <newProjectName> [-t tag]", "Clones an existing project along with all settings and configurations", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(2)
		cmd.Flags().StringVarP(&Flags.Path, "tag", "t", "", "An optional tag for the new project, only valid value is production")
		cmd.PreRunE = shared.ProjectPreRun
	})

	parent.AddCommand(project)
}
