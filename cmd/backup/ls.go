package backup

import (
	"log/slog"

	"github.com/huandu/go-sqlbuilder"
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

	db, err := util.OpenSqliteDatabaseIfExists(databaseName)
	if err != nil {
		return err
	}
	defer db.Close()

	sb := sqlbuilder.NewSelectBuilder().Select("rowid", "date", "comment").From("backups")
	selStr, selArgs := sb.Build()

	slog.Info("Querying database backups.", "log_code", "407f8430", "databaseName", databaseName, "SQL", selStr, "Args", selArgs)

	res, err := db.Query(selStr, selArgs...)
	if err != nil {
		slog.Warn("Could not query database backups.", "log_code", "5cd94c57", "databaseName", databaseName)
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

	slog.Info("Successfully ran backup ls command.", "log_code", "042af672")
	return nil
}
