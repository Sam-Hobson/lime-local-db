package db

import (
	"github.com/sam-hobson/internal/database"
	"github.com/sam-hobson/internal/types"
	"github.com/spf13/cobra"
)

// TODO: Flesh out use/examples documentation.
func newDbCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "new [DatabaseName] [DatabaseColumns]",
		Short:   "Create a new database.",
		Example: "limedb db new pets PN:TEXT:name :INT:age{0} :TEXT:gender{F}",
		Args:    cobra.MinimumNArgs(2),

		RunE: runNewDbCommand,
	}

	return cmd
}

func runNewDbCommand(cmd *cobra.Command, args []string) error {
	name := args[0]
	colStrings := args[1:]

	columns := make([]*types.Column, len(colStrings))

	for i, colStr := range colStrings {
		column, err := types.ParseColumnString(colStr)

		if err != nil {
			return err
		}

		columns[i] = column
	}

	return database.CreateDatabase(name, columns)
}
