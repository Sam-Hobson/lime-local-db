package rmentriesall

import "github.com/spf13/cobra"


func NewCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use: "rm-entries-all",
        Short: "Remove all entries from a database",
        Example: "limedb --db pets rm-entries-all",
        Args: cobra.ExactArgs(0),

        RunE: run,
    }

    return cmd
}


func run(cmd *cobra.Command, args []string) error {

    return nil
}
