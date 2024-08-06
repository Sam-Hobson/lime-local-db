package backup

import (
	"github.com/spf13/cobra"
)

// TODO: Better documentation
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "backup [Subcommand]",
		Short:   "Backup databases",
		Example: "limedb --db pets backup new",
	}

	cmd.AddCommand(
		newBackupCommand(),
		lsBackupCommand(),
		rmBackupCommand(),
		restoreBackupCommand(),
	)

	return cmd
}
