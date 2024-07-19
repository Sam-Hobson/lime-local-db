package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "limedb",
	Short: "An application for interacting with a simple database",
	Long:  "TODO: This",
}

func ProcessArgs() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(SetupCommand)
	rootCmd.AddCommand(NewDbCommand)
	rootCmd.AddCommand(RmDbCommand)
}
