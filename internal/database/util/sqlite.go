package util

import (
	"database/sql"
	"path/filepath"

	"github.com/go-errors/errors"
	"github.com/huandu/go-sqlbuilder"
	"github.com/sam-hobson/internal/types"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/viper"
)

func CreateSqliteTable(tableName string, columns []*types.Column) (string, []interface{}) {
	util.Log("19529bb3").Info("Creating sqlite to create table.", "Table name", tableName, "Columns", columns)

	ctb := sqlbuilder.NewCreateTableBuilder()
	ctb.CreateTable(tableName).IfNotExists()

	for _, col := range columns {
		opts := make([]string, 10)
		opts = append(opts, col.ColName)
		opts = append(opts, col.DataType.String())

		if col.DefaultVal != "" {
			opts = append(opts, "DEFAULT "+col.DefaultVal)
		}

		if col.NotNull {
			opts = append(opts, "NOT NULL")
		}

		if col.PrimaryKey {
			opts = append(opts, "PRIMARY KEY")
		}

		// if col.ForeignKey {
		// 	opts = append(opts, "FOREIGN KEY")
		// }

		ctb.Define(opts...)
	}

	return ctb.Build()
}

func InsertIntoSqliteTable(tableName string, entries map[string]string) (string, []interface{}) {
	util.Log("9391c009").Info("Creating sqlite to insert into table.", "Table name", tableName, "Entries", entries)

	ib := sqlbuilder.NewInsertBuilder()
	ib.InsertInto(tableName)

	keys := make([]string, len(entries))
	values := make([]interface{}, len(entries))
	var i = 0

	for key, value := range entries {
		keys[i] = key
		values[i] = value
		i++
	}

	ib.Cols(keys...)
	ib.Values(values...)

	return ib.Build()
}

func OpenSqliteDatabaseIfExists(databaseName string) (*sql.DB, error) {
    exists, err := SqliteDatabaseExists(databaseName)

	if err != nil {
		util.Log("94f1ece2").Error("Cannot backup database as it does not exist.")
        return nil, err
	}
    if !exists {
		util.Log("cbd713ec").Error("Cannot backup database as it does not exist.")
		return nil, errors.Errorf("Cannot backup database as it does not exist")
    }

    return OpenSqliteDatabase(databaseName)
}

func SqliteDatabaseExists(databaseName string) (bool, error) {
	fileName := databaseName + ".db"
	relFs := util.NewRelativeFsManager(viper.GetString("limedb_home"))
	util.Log("e75f8412").Info("Checking if database exists.", "Db name", databaseName)
	return relFs.FileExists("stores", fileName)
}

func OpenSqliteDatabase(databaseName string) (*sql.DB, error) {
	fileName := databaseName + ".db"
	dbPath := filepath.Join(viper.GetString("limedb_home"), "stores", fileName)
	util.Log("34503562").Info("Opening database file.", "Db path", dbPath)
	return sql.Open("sqlite3", dbPath)
}

// TODO: This should be refactored into a struct
func AllExistingDatabaseNames() ([]string, error) {
	relFs := util.NewRelativeFsManager(viper.GetString("limedb_home"))
	files, err := relFs.ReadDir("stores")

	if err != nil {
		return nil, err
	}

	dbNames := make([]string, len(files))

	for i, file := range files {
		dbNames[i] = file.Name()
	}

	return dbNames, nil
}
