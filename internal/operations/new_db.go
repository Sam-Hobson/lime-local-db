package operations

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/go-errors/errors"
	"github.com/huandu/go-sqlbuilder"
	_ "github.com/mattn/go-sqlite3"

	conf "github.com/sam-hobson/internal/config"
)

const StoresRelDir = "stores"

func NewDb(dbName string, columns []*Column) error {
	slog.Info("Beginning new-db operation.",
		"log_code", "26cd37c1",
		"db-name", dbName,
		"Columns", columns)

	fileName := dbName

	if !strings.HasSuffix(dbName, ".db") {
		fileName += ".db"
	}

	if err := createDbFile(StoresRelDir, fileName); err != nil {
		return err
	}

	dbPath := conf.FullPath(StoresRelDir, fileName)
	db, err := sql.Open("sqlite3", dbPath)
	defer db.Close()

	if err != nil {
		slog.Error("Could not open sqlite database.", "log_code", "9494fc60", "Db_path", dbPath)
		return err
	}

	ctb := sqlbuilder.NewCreateTableBuilder()
	ctb.CreateTable(dbName).IfNotExists()

	for _, col := range columns {
		opts := make([]string, 10)
		opts = append(opts, col.ColName)

		if col.DataType == ColumnVarCharDataType {
			opts = append(opts, fmt.Sprintf("%s(%d)", col.DataType.String(), col.VarCharLength))
		} else {
			opts = append(opts, col.DataType.String())
		}

		if col.NotNull {
			opts = append(opts, "NOT NULL")
		}

		if col.PrimaryKey {
			opts = append(opts, "PRIMARY KEY")
		}

		if col.ForeignKey {
			opts = append(opts, "FOREIGN KEY")
		}

		ctb.Define(opts...)
	}

	createTableStr, args := ctb.Build()

	slog.Info("Creating table with SQL command.", "log_code", "0cb6a54d", "SQL", createTableStr)

	_, err = db.Exec(createTableStr, args...)

	if err != nil {
		slog.Error("Failed executing create table command.", "log_code", "fed4e102", "SQL", createTableStr)
		return err
	}

	slog.Info("Successfully created a new-db.", "log_code", "7bf9634b", "Db_path", dbPath)
	return err
}

func createDbFile(fileDir, fileName string) error {
	dbPath := conf.FullPath(fileDir, fileName)

	exists, err := conf.FileExists(fileDir, fileName)

	if err != nil {
		return err
	}

	if exists {
		slog.Error("Cannot create a new database as it already exists.", "log_code", "6c95edf6", "path", dbPath)
		return errors.Errorf("Cannot create a new database as it already exists.")
	}

	err = conf.CreateFile(fileDir, fileName)

	return err
}
