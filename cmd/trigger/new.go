package trigger

import (
	"github.com/go-errors/errors"
	"github.com/sam-hobson/internal/database"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/cobra"
)

func newTriggerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "new",
		Short:   "Create/add a new trigger.",
		Example: "limedb trigger new -f mytrigger.sqlite",
		Args:    cobra.MaximumNArgs(0),

		RunE: runNewTriggerCommand,
	}

	cmd.Flags().StringP("from-file", "f", "", "Add a trigger from the file within the <LIMEDB HOME>/triggers/<db name> directory")

	return cmd
}

func runNewTriggerCommand(cmd *cobra.Command, args []string) error {
	fileName := util.PanicIfErr(cmd.Flags().GetString("from-file"))
	databaseName := state.ApplicationState().GetSelectedDb()

	if databaseName == "" {
		util.Log("5a4cfcc4").Error("Cannot add trigger if database is not selected.")
		return errors.Errorf("Cannot add trigger if database is not selected")
	}

	if fileName != "" {
		return database.NewTriggerFromFile(databaseName, fileName)
	}

	return nil
}
