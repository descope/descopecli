package project

import (
	"github.com/descope/descopecli/shared"
	"github.com/spf13/cobra"
)

var Flags struct {
	Path          string
	Debug         bool
	Tag           string
	SecretsInput  string
	SecretsOutput string
	FailureOutput string
	Force         bool
}

func AddCommands(parent *cobra.Command, group *cobra.Group) {
	project := shared.MakeGroupCommand(group, "project", "Commands for managing project environments")

	shared.AddCommand(project, Export, "export <projectId> [-p path]", "Export all project settings and configurations", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
		cmd.Flags().StringVarP(&Flags.Path, "path", "p", "", "The path to write the project data into")
		cmd.PreRunE = shared.ProjectPreRun
	})

	shared.AddCommand(project, Import, "import <projectId> [-p path] [-s secrets]", "Import settings and configurations and apply them to a project", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
		cmd.Flags().StringVarP(&Flags.Path, "path", "p", "", "The path to read the project data from")
		cmd.Flags().StringVar(&Flags.SecretsInput, "secrets-input", "", "The path to a JSON file with required secrets")
		cmd.PreRunE = shared.ProjectPreRun
	})

	shared.AddCommand(project, Validate, "validate <projectId> [-p path] [-s secrets]", "Validate project settings and configurations before import", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
		cmd.Flags().StringVarP(&Flags.Path, "path", "p", "", "The path to read the project data from")
		cmd.Flags().StringVar(&Flags.SecretsInput, "secrets-input", "", "The path to a JSON file with required secrets")
		cmd.Flags().StringVar(&Flags.FailureOutput, "failure-output", "", "The path to write a list of failures in case validation fails")
		cmd.Flags().StringVar(&Flags.SecretsOutput, "secrets-output", "", "The path to write any missing secrets in case validation fails")
		cmd.PreRunE = shared.ProjectPreRun
	})

	shared.AddCommand(project, Clone, "clone <existingProjectId> <newProjectName> [-t tag]", "Clones an existing project along with all settings and configurations", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(2)
		cmd.Flags().StringVarP(&Flags.Path, "tag", "t", "", "An optional tag for the new project, only valid value is production")
		cmd.PreRunE = shared.ProjectPreRun
	})

	shared.AddCommand(project, Delete, "delete <projectId> [-f]", "Deletes an existing project", func(cmd *cobra.Command) {
		cmd.Args = cobra.ExactArgs(1)
		cmd.Flags().BoolVarP(&Flags.Force, "force", "f", false, "Skips the prompt the deletes the project immediately")
		cmd.PreRunE = shared.ProjectPreRun
	})

	parent.AddCommand(project)
}
