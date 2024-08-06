package db

import "github.com/spf13/cobra"

// TODO: Flesh out use/examples documentation.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "db [Subcommand]",
		Short:     "Operate on databases",
		Example:   "limedb db new pets PN:TEXT:name :INT:age{0} :TEXT:gender{F}",
	}

	cmd.AddCommand(
        newDbCommand(),
        rmDbCommand(),
    )

	return cmd
}
