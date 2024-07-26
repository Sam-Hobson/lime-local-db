package rmdb

import (
	"log/slog"

	"github.com/sam-hobson/internal/database"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm-db [Database name]",
		Short:   "Delete a database",
		Example: "limedb rm-db petdb",

		Args: cobra.ExactArgs(1),
		RunE: run,
	}

	cmd.Flags().Bool("confirm", false, "Confirm that you want to take the current risky action")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	confirmed := util.PanicIfErr(cmd.Flags().GetBool("confirm"))

	if !confirmed {
		slog.Warn("rm-db rejected. Operation was not confirmed.", "log_code", "a28ce317")
		cmd.PrintErrln("Command rejected. Please use the --confirm flag if you are sure you want to proceed.")
		return nil
	}

	return database.RemoveDatabase(args[0])
}
