package newdb

import (
	"github.com/sam-hobson/internal/database"
	"github.com/spf13/cobra"
)

// TODO: Flesh out use/examples documentation.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "new-db DatabaseName DatabaseColumns",
		Short:   "Create a new database.",
		Example: "limedb new-db pets PN:TEXT:name :INT:age{0} :TEXT:gender{F}",
		Args:    cobra.MinimumNArgs(2),

		RunE: run,
	}

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	name := args[0]
	colStrings := args[1:]

	columns := make([]*database.Column, len(colStrings))

	for i, colStr := range colStrings {
		column, err := database.ParseColumnString(colStr)

		if err != nil {
			return err
		}

		columns[i] = column
	}

	return database.CreateDatabase(name, columns)
}
