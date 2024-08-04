package database

import (
	"database/sql"
	"log/slog"

	"github.com/go-errors/errors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sam-hobson/internal/types"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/viper"
)

func CreateDatabase(databaseName string, columns []*types.Column) error {
	slog.Info("Beginning new-db operation.",
		"Log code", "26cd37c1",
		"Db name", databaseName,
		"Columns", columns)

	fileName := databaseName + ".db"
	relFs := util.NewRelativeFsManager(viper.GetString("limedb_home"))

	if exists, err := relFs.FileExists("stores", fileName); err != nil {
		return err
	} else if exists {
		slog.Error("Cannot create a new database as it already exists.", "Log code", "6c95edf6", "Database name", databaseName)
		return errors.Errorf("Cannot create a new database as it already exists")
	}

	if err := relFs.CreateFile("stores", fileName); err != nil {
		return err
	}

	dbPath := relFs.FullPath("stores", fileName)
	db, err := sql.Open("sqlite3", dbPath)
	defer db.Close()

	if err != nil {
		slog.Error("Could not open sqlite database.", "Log code", "9494fc60", "Db path", dbPath)
		return err
	}

	createTableStr, args := util.CreateSqliteTable(databaseName, columns)

	slog.Info("Creating table with SQL command.", "Log code", "0cb6a54d", "SQL", createTableStr, "Args", args)

	if _, err = db.Exec(createTableStr, args...); err != nil {
		slog.Error("Failed executing create table command.", "Log code", "fed4e102", "SQL", createTableStr)
		return err
	}

	slog.Info("Successfully created a new database.", "Log code", "7bf9634b", "Db path", dbPath)
	return nil
}
