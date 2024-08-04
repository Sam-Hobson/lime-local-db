package backup

import (
	"github.com/go-errors/errors"
	"github.com/huandu/go-sqlbuilder"
	"github.com/sam-hobson/internal/database"
	dbutil "github.com/sam-hobson/internal/database/util"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/cobra"
)

func createNewBackupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "new",
		Short:     "Backup a database",
		Example:   "limedb backup new",
		Args:      cobra.ExactArgs(0),
		ValidArgs: dbNames(),

		RunE: runNewBackupCommand,
	}

	cmd.Flags().StringP("message", "m", "", "Add a message/note associated with the backup")

	return cmd
}

func runNewBackupCommand(cmd *cobra.Command, args []string) error {
	databaseName := state.ApplicationState().GetSelectedDb()

	if databaseName == "" {
		util.Log("d6eda884").Error("Cannot list backups if no database selected.")
		return errors.Errorf("Cannot list backups if no database selected")
	}

	if lastId, err := database.BackupDatabase(databaseName, util.PanicIfErr(cmd.Flags().GetString("message"))); err != nil {
		return err
	} else {
		cond := sqlbuilder.NewCond()
		where := sqlbuilder.NewWhereClause()
		where.AddWhereExpr(cond.Args, cond.Equal("rowid", lastId))
		printBackupRowsWhere(cmd, databaseName, where)
	}

	return nil
}

func dbNames() []string {
	if names, err := dbutil.AllExistingDatabaseNames(); err != nil {
		return []string{}
	} else {
		return names
	}
}
