package trigger

import (
	"github.com/huandu/go-sqlbuilder"
	dbutil "github.com/sam-hobson/internal/database/util"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/cobra"
)

// TODO: Flesh out use/examples documentation.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "trigger [Subcommand]",
		Short:             "Operate on database triggers",
		Example:           "limedb trigger",
		PersistentPreRunE: persistentPreRun,
	}

	cmd.AddCommand(
		newTriggerCommand(),
		templateTriggerCommand(),
		lsTriggerCommand(),
	)

	return cmd
}

func persistentPreRun(cmd *cobra.Command, _ []string) error {
	if err := syncTriggersTable(); err != nil {
		return err
	}

	return nil
}

func syncTriggersTable() error {
	util.Log("555cc716").Info("Checking sqlite_master for triggers.")

	databaseName := state.ApplicationState().GetSelectedDb()
	if databaseName == "" {
		return nil
	}

	cond := sqlbuilder.NewCond()
	sql, args := dbutil.EntriesInTableWhereSql("sqlite_master", []string{"rowid", "name"}, cond.Args, cond.Equal("type", "trigger"))

	db, err := dbutil.OpenSqliteDatabaseIfExists(databaseName)
	if err != nil {
		return err
	}

	util.Log("161517c1").Info("Querying sqlite master with SQL.", "Database name", databaseName, "SQL", sql, "Args", args)

	rows, err := db.Query(sql, args...)
	if err != nil {
		util.Log("fed59b67").Error("Could not query sqlite_master table.", "Database name", databaseName)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var rowid int
		var name string
		if err := rows.Scan(&rowid, &name); err != nil {
			util.Log("44b3ba8d").Error("Error occurred while reading from sqlite_master.", "Database name", databaseName)
			return err
		}
	}

	return nil
}
