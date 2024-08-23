package database

import (
	"github.com/go-errors/errors"
	_ "github.com/mattn/go-sqlite3"
	dbutil "github.com/sam-hobson/internal/database/util"
	"github.com/sam-hobson/internal/types"
	"github.com/sam-hobson/internal/util"
)

var backupColumns = []*types.Column{
	{
		Name:     "date",
		DataType: types.ColumnTextDataType,
		NotNull:  true,
	},
	{
		Name:     "backup_name",
		DataType: types.ColumnTextDataType,
		NotNull:  true,
	},
	{
		Name:     "comment",
		DataType: types.ColumnTextDataType,
	},
}

var triggerColumns = []*types.Column{
	{
		Name:       "sqlite_master_rowid",
		DataType:   types.ColumnIntDataType,
		NotNull:    true,
		PrimaryKey: true,
		ForeignKey: &types.ForeignKey{
			Table: "sqlite_master",
			Col:   "rowid",
		},
	},
	{
		Name:     "date_created",
		DataType: types.ColumnTextDataType,
		NotNull:  true,
	},
	{
		Name:     "trigger_type",
		DataType: types.ColumnTextDataType,
	},
	{
		Name:     "comment",
		DataType: types.ColumnTextDataType,
	},
}

func CreateDatabase(databaseName string, columns []*types.Column) error {
	util.Log("26cd37c1").Info("Beginning db new operation.", "Database name", databaseName, "Columns", columns)

	if exists, err := dbutil.SqliteDatabaseExists(databaseName); err != nil {
		return err
	} else if exists {
		util.Log("6c95edf6").Error("Cannot create a new database as it already exists.", "Database name", databaseName)
		return errors.Errorf("Cannot create a new database as it already exists")
	}

	db, err := dbutil.OpenSqliteDatabase(databaseName)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create data table for db
	createTableStr, args := dbutil.CreateTableSql(databaseName, columns)
	util.Log("0cb6a54d").Info("Creating table with SQL command.", "SQL", createTableStr, "Args", args)

	if _, err = db.Exec(createTableStr, args...); err != nil {
		util.Log("fed4e102").Error("Failed executing create table command.", "SQL", createTableStr, "Args", args)
		return err
	}

	util.Log("7bf9634b").Info("Successfully created a new database.", "Database name", databaseName)

	// Create triggers table
	createTriggersStr, triggersArgs := dbutil.CreateTableSql("triggers", triggerColumns)
	util.Log("eadfbde6").Info("Creating trigger table with SQL.", "SQL", createTriggersStr, "Args", triggersArgs)

	if _, err := db.Exec(createTriggersStr, triggersArgs...); err != nil {
		util.Log("605f11a0").Error("Failed creating trigger table.", "SQL", createTriggersStr, "Args", triggersArgs)
		dbutil.RemoveSqliteDatabase(databaseName)
		return err
	}

	// Create accompanying persistent database
	if err := CreatePersistentDatabase(databaseName); err != nil {
		dbutil.RemoveSqliteDatabase(databaseName)
		return err
	}

	return nil
}

func CreatePersistentDatabase(databaseName string) error {
	persistentDatabaseName := dbutil.PersistentDatabaseName(databaseName)
	util.Log("3e55ef45").Info("Create persistent database operation.", "Database name", databaseName, "Persistent database name", persistentDatabaseName)

	if exists, err := dbutil.SqliteDatabaseExists(persistentDatabaseName); err != nil {
		return err
	} else if exists {
		util.Log("70de0695").Error("Cannot create a new database as a persistent database already exists.", "Persistent database name", persistentDatabaseName)
		return errors.Errorf("Cannot create a new database as a persistent database already exists")
	}

	db, err := dbutil.OpenSqliteDatabase(persistentDatabaseName)
	if err != nil {
		util.Log("9494fc60").Error("Could not open sqlite database.", "Persistent database name", persistentDatabaseName)
		return err
	}
	defer db.Close()

	createTableStr, createTableArgs := dbutil.CreateTableSql("backups", backupColumns)
	util.Log("8bc1e038").Info("Creating backup table with SQL command.", "SQL", createTableStr, "Args", createTableArgs)

	if _, err := db.Exec(createTableStr, createTableArgs...); err != nil {
		util.Log("f7d58d42").Error("Failed executing create table command.", "SQL", createTableStr)
		return err
	}

	util.Log("91750756").Info("Successfully created a backup table.")

	return nil
}
