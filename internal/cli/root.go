package cli

import (
	"github.com/sam-hobson/internal/operations"
	"github.com/spf13/cobra"
)

var selectedDb = ""
var selectedCol = ""

var rootCmd = &cobra.Command{
	Use:   "limedb",
	Short: "An application for interacting with a simple database",
	Long:  "TODO: This",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        if cmd.PersistentFlags().Lookup("db").Changed {
            operations.SelectDb(selectedDb)
        }

        if selectedCol != "" {
            operations.SelectCol(selectedCol)
        }

        return nil
	},
}

func ProcessArgs() error {
	return rootCmd.Execute()
}

func init() {
	globalFlags := rootCmd.PersistentFlags()

	globalFlags.StringVar(&selectedDb, "db", "", "Selected database to use for operations")
	globalFlags.StringVar(&selectedCol, "col", "", "Selected column to use for operations")

	rootCmd.AddCommand(SetupCommand)

	rootCmd.AddCommand(NewDbCommand)
	rootCmd.AddCommand(RmDbCommand)
}
