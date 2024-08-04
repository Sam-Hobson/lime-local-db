package backup

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "backup [Subcommand]",
		Short:   "Backup a database",
		Example: "limedb backup petdb",
	}

	cmd.AddCommand(
        createNewBackupCommand(),
        lsBackupCommand(),
    )

	return cmd
}
