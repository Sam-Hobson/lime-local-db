package rmdb

import (
	"github.com/sam-hobson/internal/database"
	dbutil "github.com/sam-hobson/internal/database/util"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "rm-db [Database name]",
		Short:     "Delete a database",
		Example:   "limedb rm-db petdb",
		Args:      cobra.ExactArgs(1),
		ValidArgs: dbNames(),

		RunE: run,
	}

	cmd.Flags().Bool("confirm", false, "Confirm that you want to take the current risky action")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	if !util.PanicIfErr(cmd.Flags().GetBool("confirm")) {
		util.Log("a28ce317").Warn("rm-db rejected. Operation was not confirmed.")
		cmd.PrintErrln("Command rejected. Please use the --confirm flag if you are sure you want to proceed.")
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
