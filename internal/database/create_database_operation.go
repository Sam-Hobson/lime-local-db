package database

import (
	"database/sql"
	"log/slog"
	"strings"

	"github.com/go-errors/errors"
	"github.com/huandu/go-sqlbuilder"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/viper"
)

func CreateDatabase(databaseName string, columns []*Column) error {
	slog.Info("Beginning new-db operation.",
		"log_code", "26cd37c1",
		"db-name", databaseName,
		"Columns", columns)

	fileName := databaseName

	if !strings.HasSuffix(databaseName, ".db") {
		fileName += ".db"
	}

	relFs := util.NewRelativeFsManager(viper.GetString("limedbHome"))

	if exists, err := relFs.FileExists("stores", fileName); err != nil {
		return err
	} else if exists {
		slog.Error("Cannot create a new database as it already exists.", "log_code", "6c95edf6", "database_name", databaseName)
		return errors.Errorf("Cannot create a new database as it already exists")
	}

	if err := relFs.CreateFile("stores", fileName); err != nil {
		return err
	}

	dbPath := relFs.FullPath("stores", fileName)
	db, err := sql.Open("sqlite3", dbPath)
	defer db.Close()

	if err != nil {
		slog.Error("Could not open sqlite database.", "log_code", "9494fc60", "Db_path", dbPath)
		return err
	}

	ctb := sqlbuilder.NewCreateTableBuilder()
	ctb.CreateTable(databaseName).IfNotExists()

	for _, col := range columns {
		opts := make([]string, 10)
		opts = append(opts, col.ColName)
		opts = append(opts, col.DataType.String())

        if col.DefaultVal != "" {
            opts = append(opts, "DEFAULT " + col.DefaultVal)
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

	createTableStr, args := ctb.Build()

	slog.Info("Creating table with SQL command.", "log_code", "0cb6a54d", "SQL", createTableStr)

	if _, err = db.Exec(createTableStr, args...); err != nil {
		slog.Error("Failed executing create table command.", "log_code", "fed4e102", "SQL", createTableStr)
		return err
	}

	slog.Info("Successfully created a new database.", "log_code", "7bf9634b", "Db_path", dbPath)
	return err

}
