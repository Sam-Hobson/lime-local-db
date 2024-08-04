package backup

import (
	"github.com/huandu/go-sqlbuilder"
	dbutil "github.com/sam-hobson/internal/database/util"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/cobra"
)

func lsBackupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "ls [Database name]",
		Short:     "List backups for a given database",
		Example:   "limedb backup ls pets",
		Args:      cobra.ExactArgs(1),
		ValidArgs: dbNames(),

		RunE: runLsBackupCommand,
	}

	return cmd
}

func runLsBackupCommand(cmd *cobra.Command, args []string) error {
	databaseName := args[0]

	db, err := dbutil.OpenSqliteDatabaseIfExists(databaseName)
	if err != nil {
		return err
	}
	defer db.Close()

	sb := sqlbuilder.NewSelectBuilder().Select("rowid", "date", "comment").From("backups")
	selStr, selArgs := sb.Build()

	util.Log("6fac1254").Info("Querying database backups.", "Database name", databaseName, "SQL", selStr, "Args", selArgs)

	res, err := db.Query(selStr, selArgs...)
	if err != nil {
		util.Log("5cd94c57").Warn("Could not query database backups.", "Database name", databaseName)
		return err
	}
	defer res.Close()

	for res.Next() {
		var rowid int
		var date string
		var comment string

		res.Scan(&rowid, &date, &comment)

		cmd.Printf("%d  %s  \"%s\"\n", rowid, date, comment)
	}

	util.Log("5dbf67eb").Info("Successfully ran backup ls command.")
	return nil
}
