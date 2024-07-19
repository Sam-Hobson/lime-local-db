package cli

import (
	"log/slog"

	op "github.com/sam-hobson/internal/operations"
	"github.com/spf13/cobra"
)

var RmDbCommand = &cobra.Command{
	Use:     "rm-db [Database name]",
	Example: "limedb rm-db petdb",
	Short:   "Delete a database",
	Args:    cobra.ExactArgs(1),

	RunE: processRmDbCmd,
}

func processRmDbCmd(cmd *cobra.Command, args []string) error {
	confirmed := cmd.Flags().Changed("confirm")

	if !confirmed {
		slog.Warn("rm-db rejected. Operation was not confirmed.", "log_code", "a28ce317")
		cmd.PrintErrln("Command rejected. Please use the --confirm flag if you are sure you want to proceed.")
		return nil
	}

	return op.RmDb(args[0])
}

func init() {
	RmDbCommand.Flags().Bool("confirm", false, "Confirm that you want to take the current risky action")
}
