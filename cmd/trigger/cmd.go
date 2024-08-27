package trigger

import (
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/sam-hobson/internal/database"
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

	db, err := dbutil.OpenSqliteDatabaseIfExists(databaseName)
	if err != nil {
		return err
	}
	defer db.Close()

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("l.rowid")
	sb.From("sqlite_master l")
	sb.JoinWithOption(sqlbuilder.LeftJoin, "triggers r", "l.rowid = r.sqlite_master_rowid")
	sb.Where("l.type = "+sb.Var("trigger"), "r.sqlite_master_rowid IS NULL")

	sqliteMasterStr, sqliteMasterArgs := sb.Build()

	util.Log("a5363de7").Info("Querying sqlite_master with following SQL.", "Database name", databaseName, "SQL", sqliteMasterStr, "Args", sqliteMasterArgs)

	masterRowidRows, err := db.Query(sqliteMasterStr, sqliteMasterArgs...)
	if err != nil {
		util.Log("fed59b67").Error("Could not query sqlite_master table for rowids.", "Database name", databaseName)
		return err
	}
	defer masterRowidRows.Close()

	masterRowIds, err := dbutil.RowsIntoSlice[int](masterRowidRows)
	if err != nil {
		return err
	}

	now := time.Now().Format(time.RFC3339)
	triggers := make([]interface{}, len(masterRowIds))
	for i, id := range masterRowIds {
		triggers[i] = &database.Trigger{
			SqliteMasterRowid: id,
			Date:              now,
		}
	}

	insertStr, args := database.TriggerStruct.InsertInto("triggers", triggers...).Build()
	util.Log("2100aa9e").Info("Inserting into triggers with SQL.", "Database name", databaseName, "SQL", insertStr, "Args", args)
	if _, err := db.Exec(insertStr, args...); err != nil {
		util.Log("b12b7f7d").Error("Failed ")
	}

	// sqliteMasterCond := sqlbuilder.NewCond()
	// selectSqliteMasterStr, sqliteMasterArgs := dbutil.EntriesInTableWhereSql("sqlite_master", []string{"rowid"}, sqliteMasterCond.Args, sqliteMasterCond.Equal("type", "trigger"))
	//
	// util.Log("161517c1").Info("Querying sqlite master with SQL.", "Database name", databaseName, "SQL", selectSqliteMasterStr, "Args", sqliteMasterArgs)
	//
	// sqliteMasterRows, err := db.Query(selectSqliteMasterStr, sqliteMasterArgs...)
	// if err != nil {
	// 	util.Log("fed59b67").Error("Could not query sqlite_master table.", "Database name", databaseName)
	// 	return err
	// }
	// defer sqliteMasterRows.Close()
	//
	// rowids, err := dbutil.RowsIntoSlice[int](sqliteMasterRows)
	// if err != nil {
	// 	util.Log("f6f64c29").Error("Could not reading sqlite_master rowids into slice.", "Database name", databaseName)
	// 	return err
	// }
	//
	// triggersCond := sqlbuilder.NewCond()
	// triggersStr, triggersArgs := dbutil.EntriesInTableWhereSql("triggers", )

	return nil
}
