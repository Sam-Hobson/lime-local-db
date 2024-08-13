package rmentriesall

import (
	"github.com/go-errors/errors"
	"github.com/huandu/go-sqlbuilder"
	"github.com/sam-hobson/internal/database"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm-entries-all",
		Short:   "Remove all entries from a database",
		Example: "limedb --db pets rm-entries-all",
		Args:    cobra.ExactArgs(0),

		RunE: run,
	}

	cmd.Flags().Bool("confirm", false, "Confirm that you want to take the current risky action")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	if !util.PanicIfErr(cmd.Flags().GetBool("confirm")) {
		util.Log("3912ef9e").Warn("rm-entries-all rejected. Operation was not confirmed.")
		cmd.PrintErrln("Command rejected. Please use the --confirm flag if you are sure you want to proceed.")
		return nil
	}

	databaseName := state.ApplicationState().GetSelectedDb()
	if databaseName == "" {
		util.Log("f3878ee8").Error("Cannot remove entries if not database is specified.")
		return errors.Errorf("Cannot remove entries if not database is specified")
	}

	if rowsEffected, err := database.RemoveEntries(databaseName, sqlbuilder.NewWhereClause()); err != nil {
		return err
	} else {
		cmd.Printf("%d rows affected\n", rowsEffected)
		return nil
	}
}
