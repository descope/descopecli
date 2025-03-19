package shared

import (
	"os"

	"github.com/spf13/cobra"
)

func AddCommand(parent *cobra.Command, action func([]string) error, use string, help string, setup func(*cobra.Command)) {
	cmd := &cobra.Command{
		Use:                   use,
		Short:                 help,
		DisableFlagsInUseLine: true,
		PreRunE:               DefaultPreRun,
		Run: func(_ *cobra.Command, args []string) {
			err := action(args)
			ExitWithStatus(err)
			if err != nil {
				os.Exit(1)
			}
		},
	}
	cmd.Args = cobra.ExactArgs(0)
	setup(cmd)
	cmd.Flags().BoolVarP(&Flags.Json, "json", "j", false, "use JSON output format")
	hideHelpFlag(cmd)
	parent.AddCommand(cmd)
}

func MakeGroupCommand(parentGroup *cobra.Group, use string, help string) *cobra.Command {
	groupID := ""
	if parentGroup != nil {
		groupID = parentGroup.ID
	}
	cmd := &cobra.Command{
		GroupID: groupID,
		Use:     use,
		Short:   help,
	}
	hideHelpFlag(cmd)
	return cmd
}

func hideHelpFlag(cmd *cobra.Command) {
	cmd.InitDefaultHelpFlag()
	_ = cmd.Flags().MarkHidden("help")
	cmd.Flags().SortFlags = false
}
