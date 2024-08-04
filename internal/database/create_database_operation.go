package database

import (
	"database/sql"

	"github.com/go-errors/errors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sam-hobson/internal/types"
	"github.com/sam-hobson/internal/util"
	dbutil "github.com/sam-hobson/internal/database/util"
	"github.com/spf13/viper"
)

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
	return nil
}
