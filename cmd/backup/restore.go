package backup

import (
	"strconv"

	"github.com/go-errors/errors"
	"github.com/sam-hobson/internal/database"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/cobra"
)

func restoreBackupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "restore [Backup id]",
		Short:     "Restore a database",
		Example:   "limedb backup restore [Backup id]",
		Args:      cobra.ExactArgs(1),

		RunE: runRestoreBackupCommand,
	}

	cmd.Flags().Bool("confirm", false, "Confirm that you want to take the current risky action")

	return cmd
}

func runRestoreBackupCommand(cmd *cobra.Command, args []string) error {
	if !util.PanicIfErr(cmd.Flags().GetBool("confirm")) {
		util.Log("53f56171").Warn("backup restore rejected. Operation was not confirmed.")
		cmd.PrintErrln("Command rejected. Please use the --confirm flag if you are sure you want to proceed")
		return nil
	}

	databaseName := state.ApplicationState().GetSelectedDb()

	if databaseName == "" {
		util.Log("54cd0827").Error("Cannot restore from backups if no database selected.")
		return errors.Errorf("Cannot restore from backups if no database selected")
	}

    database.RestoreFromBackup(databaseName, util.PanicIfErr(strconv.Atoi(args[0])))

	return nil
}
