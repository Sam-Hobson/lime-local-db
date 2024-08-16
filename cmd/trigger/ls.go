package trigger

import (
	"fmt"

	"github.com/go-errors/errors"
	dbutil "github.com/sam-hobson/internal/database/util"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/cobra"
)

func lsTriggerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Short:   "Prints the triggers on the selected database",
		Example: "limedb trigger ls",
		Args:    cobra.NoArgs,

		RunE: runLsTriggerCommand,
	}

	// TODO: Implement this
	cmd.Flags().Bool("all", false, "Show all details about triggers")

	return cmd
}

func runLsTriggerCommand(cmd *cobra.Command, args []string) error {
	databaseName := state.ApplicationState().GetSelectedDb()
	util.Log("bd1477f3").Info("Listing all database triggers.", "Database name", databaseName)
	if databaseName == "" {
		util.Log("1c822e6d").Error("Cannot ls triggers if database is not selected.")
		return errors.Errorf("Cannot ls triggers if database is not selected")
	}

	triggers, err := dbutil.DefinedTriggers(databaseName)
	if err != nil {
		return err
	}

	for _, trigger := range triggers {
		fmt.Println(trigger.Name)
	}

	return nil
}
