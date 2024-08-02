package rmentriesall

import (
	"log/slog"

	"github.com/huandu/go-sqlbuilder"
	"github.com/sam-hobson/internal/database"
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
		slog.Warn("rm-db rejected. Operation was not confirmed.", "log_code", "abf450ec")
		cmd.PrintErrln("Command rejected. Please use the --confirm flag if you are sure you want to proceed.")
		return nil
	}

    if rowsEffected, err := database.RemoveEntries(sqlbuilder.NewWhereClause()); err != nil {
        return err
    } else {
        cmd.Printf("%d rows effected\n", rowsEffected)
        return nil
    }
}
