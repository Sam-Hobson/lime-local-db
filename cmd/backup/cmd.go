package backup

import (
	"github.com/sam-hobson/internal/database"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "backup [Subcommand]",
		Short:   "Backup a database",
		Example: "limedb backup petdb",

		PersistentPostRun: postBackupRun,
	}

	cmd.AddCommand(
		createNewBackupCommand(),
		lsBackupCommand(),
		rmBackupCommand(),
        restoreBackupCommand(),
	)

	return cmd
}

func postBackupRun(cmd *cobra.Command, args []string) {
	// For all backup commands, remove orphan backups if there are any
	if viper.GetBool("remove_orphan_backups") {
		databaseName := state.ApplicationState().GetSelectedDb()
		if databaseName == "" {
			util.Log("906b4c63").Error("Cannot backups if no database selected.")
			return
		}
		database.RemoveOrphanBackups(databaseName)
	}
}
