package shared

import (
	"github.com/spf13/cobra"
)

func AddCommand(parent *cobra.Command, action func([]string) error, use string, help string, setup func(*cobra.Command)) {
	cmd := &cobra.Command{
		Use:                   use,
		Short:                 help,
		DisableFlagsInUseLine: true,
		PreRunE:               DefaultPreRun,
		RunE: func(_ *cobra.Command, args []string) error {
			return action(args)
		},
	}
	cmd.Args = cobra.ExactArgs(0)
	setup(cmd)
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
	cmd.Flags().Lookup("help").Hidden = true
	cmd.Flags().SortFlags = false
}
