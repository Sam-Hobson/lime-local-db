package util

import (
	"database/sql"
	"path/filepath"

	"github.com/go-errors/errors"
	"github.com/huandu/go-sqlbuilder"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/types"
	"github.com/sam-hobson/internal/util"
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
		util.Log("94f1ece2").Error("Cannot open database as it does not exist.")
		return nil, err
	}
	if !exists {
		util.Log("cbd713ec").Error("Cannot open database as it does not exist.")
		return nil, errors.Errorf("Cannot open database as it does not exist")
	}

	return OpenSqliteDatabase(databaseName)
}

func SqliteDatabaseExists(databaseName string) (bool, error) {
	fileName := databaseName + ".db"
	relFs := util.NewRelativeFsManager(state.ApplicationState().GetLimedbHome(), "stores")
	util.Log("e75f8412").Info("Checking if database exists.", "Database name", databaseName)
	return relFs.FileExists("", fileName)
}

func OpenSqliteDatabase(databaseName string) (*sql.DB, error) {
	fileName := databaseName + ".db"
	dbPath := filepath.Join(state.ApplicationState().GetLimedbHome(), "stores", fileName)
	util.Log("34503562").Info("Opening database file.", "Database path", dbPath)
	return sql.Open("sqlite3", dbPath)
}

func RemoveSqliteDatabase(databaseName string) error {
	fileName := databaseName + ".db"
	relFs := util.NewRelativeFsManager(state.ApplicationState().GetLimedbHome(), "stores")
	util.Log("cb25f7f8").Info("Removing database file.", "Database name", databaseName)
	return relFs.RmFile("", fileName)
}
