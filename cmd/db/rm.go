package db

import (
	"github.com/sam-hobson/internal/database"
	dbutil "github.com/sam-hobson/internal/database/util"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/cobra"
)

func rmDbCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "rm [Database name]",
		Short:     "Delete a database",
		Example:   "limedb db rm petdb",
		Args:      cobra.ExactArgs(1),
		ValidArgs: dbNames(),

		RunE: runRmDbCommand,
	}

	cmd.Flags().Bool("confirm", false, "Confirm that you want to take the current risky action")

	return cmd
}

func runRmDbCommand(cmd *cobra.Command, args []string) error {
	if !util.PanicIfErr(cmd.Flags().GetBool("confirm")) {
		util.Log("a28ce317").Warn("rm-db rejected. Operation was not confirmed.")
		cmd.PrintErrln("Command rejected. Please use the --confirm flag if you are sure you want to proceed")
		return nil
	}

	return database.RemoveDatabase(args[0])
}

func dbNames() []string {
	if names, err := dbutil.AllExistingDatabaseNames(); err != nil {
		return []string{}
	} else {
		return names
	}
}
