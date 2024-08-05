package database

import (
	"database/sql"

	"github.com/go-errors/errors"
	_ "github.com/mattn/go-sqlite3"
	dbutil "github.com/sam-hobson/internal/database/util"
	"github.com/sam-hobson/internal/types"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/viper"
)

var backupColumns = []*types.Column{
	{
		ColName:  "date",
		DataType: types.ColumnTextDataType,
		NotNull:  true,
	},
	{
		ColName:  "backupName",
		DataType: types.ColumnTextDataType,
		NotNull:  true,
	},
	{
		ColName:  "comment",
		DataType: types.ColumnTextDataType,
		NotNull:  false,
	},
}

func CreateDatabase(databaseName string, columns []*types.Column) error {
	util.Log("26cd37c1").Info("Beginning new-db operation.", "Db name", databaseName, "Columns", columns)
	fileName := databaseName + ".db"
	relFs := util.NewRelativeFsManager(viper.GetString("limedb_home"))

	if exists, err := relFs.FileExists("stores", fileName); err != nil {
		return err
	} else if exists {
		util.Log("6c95edf6").Error("Cannot create a new database as it already exists.", "Database name", databaseName)
		return errors.Errorf("Cannot create a new database as it already exists")
	}

	if err := relFs.CreateFile("stores", fileName); err != nil {
		return err
	}

	dbPath := relFs.FullPath("stores", fileName)
	db, err := sql.Open("sqlite3", dbPath)
	defer db.Close()

	if err != nil {
		util.Log("9494fc60").Error("Could not open sqlite database.", "Db path", dbPath)
		return err
	}

	createTableStr, args := dbutil.CreateSqliteTable(databaseName, columns)

	util.Log("0cb6a54d").Info("Creating table with SQL command.", "SQL", createTableStr, "Args", args)

	if _, err = db.Exec(createTableStr, args...); err != nil {
		util.Log("fed4e102").Error("Failed executing create table command.", "SQL", createTableStr)
		return err
	}

	util.Log("7bf9634b").Info("Successfully created a new database.", "Db path", dbPath)

    if err := CreatePersistentDatabase(databaseName); err != nil {
        relFs.RmFile("stores", fileName)
        return err
    }

	return nil
}

func CreatePersistentDatabase(databaseName string) error {
    persistentDatabaseFileName := dbutil.PersistentDatabaseName(databaseName) + ".db"
    util.Log("3e55ef45").Info("Beginning create persistent database operation.", "Database name", databaseName, "Persistent database file name", persistentDatabaseFileName)

	relFs := util.NewRelativeFsManager(viper.GetString("limedb_home"))

	if exists, err := relFs.FileExists("stores", persistentDatabaseFileName); err != nil {
		return err
	} else if exists {
		util.Log("70de0695").Error("Cannot create a new database as a persistent database already exists.", "Persistent database file name", persistentDatabaseFileName)
		return errors.Errorf("Cannot create a new database as a persistent database already exists")
	}

	if err := relFs.CreateFile("stores", persistentDatabaseFileName); err != nil {
		return err
    }

	dbPath := relFs.FullPath("stores", persistentDatabaseFileName)

	db, err := sql.Open("sqlite3", dbPath)
	defer db.Close()

	if err != nil {
		util.Log("9494fc60").Error("Could not open sqlite database.", "Db path", dbPath)
		return err
	}

	createTableStr, createTableArgs := dbutil.CreateSqliteTable("backups", backupColumns)
	util.Log("8bc1e038").Info("Creating backup table with SQL command.", "SQL", createTableStr, "Args", createTableArgs)

	if _, err := db.Exec(createTableStr, createTableArgs...); err != nil {
		util.Log("f7d58d42").Error("Failed executing create table command.", "SQL", createTableStr)
		return err
	}

	util.Log("91750756").Info("Successfully created a backup table.")

    return nil
}
