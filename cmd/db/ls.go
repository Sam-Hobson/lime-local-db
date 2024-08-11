package db

import (
	dbutil "github.com/sam-hobson/internal/database/util"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/cobra"
)

func lsDbCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "ls",
		Short:     "List databases",
		Example:   "limedb db ls",
		Args:      cobra.NoArgs,
		ValidArgs: dbNames(),

		RunE: runLsDbCommand,
	}

	// TODO: Implement this
	cmd.Flags().Bool("all", false, "Show all details about databases")

	return cmd
}

func runLsDbCommand(cmd *cobra.Command, _ []string) error {
	util.Log("7b1030a2").Info("Listing all databases in /store.")
	names, err := dbutil.AllExistingDatabaseNames()
	if err != nil {
		return nil
	}

	selectedDb := state.ApplicationState().GetSelectedDb()

	for _, name := range names {
		if selectedDb != "" {
			if name == selectedDb {
				cmd.Println(name)
			}
		} else {
			cmd.Println(name)
		}
	}

	return nil
}
