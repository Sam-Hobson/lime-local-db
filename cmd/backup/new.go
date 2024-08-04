package backup

import (
	"github.com/sam-hobson/internal/database"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/cobra"
)

func createNewBackupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "new [Database name]",
		Short:     "Backup a database",
		Example:   "limedb backup new pets",
		Args:      cobra.ExactArgs(1),
		ValidArgs: dbNames(),

		RunE: runNewBackupCommand,
	}

	cmd.Flags().StringP("message", "m", "", "Add a message/note associated with the backup")

	return cmd
}

func runNewBackupCommand(cmd *cobra.Command, args []string) error {
	return database.BackupDatabase(args[0], util.PanicIfErr(cmd.Flags().GetString("message")))
}

func dbNames() []string {
	if names, err := util.AllExistingDatabaseNames(); err != nil {
		return []string{}
	} else {
		return names
	}
}
