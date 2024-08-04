package backup

import (
	"strconv"

	"github.com/go-errors/errors"
	"github.com/sam-hobson/internal/database"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/cobra"
)

func rmBackupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "rm [Backup id]",
		Short:     "Remove a database backup",
		Example:   "limedb backup rm 1",
		Args:      cobra.ExactArgs(1),

		RunE: runRmBackupCommand,
	}

	return cmd
}

func runRmBackupCommand(cmd *cobra.Command, args []string) error {
	databaseName := state.ApplicationState().GetSelectedDb()

	if databaseName == "" {
		util.Log("d2b56915").Error("Cannot list backups if no database selected.")
		return errors.Errorf("Cannot list backups if no database selected")
	}

    return database.RemoveDatabaseBackup(databaseName, util.PanicIfErr(strconv.Atoi(args[0])))
}
