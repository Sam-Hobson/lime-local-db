package backup

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/go-errors/errors"
	"github.com/huandu/go-sqlbuilder"
	dbutil "github.com/sam-hobson/internal/database/util"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/cobra"
)

func lsBackupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "ls",
		Short:     "List backups for a selected database",
		Example:   "limedb backup ls",
		Args:      cobra.ExactArgs(0),
		ValidArgs: dbNames(),

		RunE: runLsBackupCommand,
	}

	return cmd
}

func runLsBackupCommand(cmd *cobra.Command, args []string) error {
	databaseName := state.ApplicationState().GetSelectedDb()

	if databaseName == "" {
		util.Log("7ddeeaef").Error("Cannot list backups if no database selected.")
		return errors.Errorf("Cannot list backups if no database selected")
	}

	if err := printBackupRowsWhere(cmd.OutOrStdout(), databaseName, sqlbuilder.NewWhereClause()); err != nil {
		return err
	}

	util.Log("5dbf67eb").Info("Successfully ran backup ls command.")
	return nil
}

func printBackupRowsWhere(iowriter io.Writer, databaseName string, where *sqlbuilder.WhereClause) error {
	PersistentDatabaseName := dbutil.PersistentDatabaseName(databaseName)
	db, err := dbutil.OpenSqliteDatabaseIfExists(PersistentDatabaseName)
	if err != nil {
		return err
	}
	defer db.Close()

	sb := sqlbuilder.NewSelectBuilder().Select("rowid", "date", "comment").From("backups").AddWhereClause(where)
	selStr, selArgs := sb.Build()

	util.Log("6fac1254").Info("Querying database backups.", "Database name", databaseName, "SQL", selStr, "Args", selArgs)

	res, err := db.Query(selStr, selArgs...)
	if err != nil {
		util.Log("5cd94c57").Warn("Could not query database backups.", "Database name", databaseName)
		return err
	}
	defer res.Close()

	w := tabwriter.NewWriter(iowriter, 2, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Id\tDate\tComment\n")

	for res.Next() {
		var rowid int
		var date string
		var comment string

		if err := res.Scan(&rowid, &date, &comment); err != nil {
			util.Log("d92e36ce").Warn("Error while reading backup results from database.", "Database name", databaseName)
			return errors.Errorf("Error while reading backup results from database")
		}
		fmt.Fprintf(w, "%d\t%s\t%s\n", rowid, date, comment)
	}

	w.Flush()

	return nil
}
