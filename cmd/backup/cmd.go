package backup

import (
	"github.com/sam-hobson/internal/database"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "backup [Database name]",
		Short:     "Backup a database",
		Example:   "limedb backup petdb",
		Args:      cobra.ExactArgs(1),
		ValidArgs: dbNames(),

		RunE: run,
	}

	cmd.Flags().StringP("message", "m", "", "Add a message/note associated with the backup")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	return database.BackupDatabase(args[0], util.PanicIfErr(cmd.Flags().GetString("message")))
}

func dbNames() []string {
	if names, err := util.AllExistingDatabaseNames(); err != nil {
		return []string{}
	} else {
		return names
	}
}
