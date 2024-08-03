package database

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/go-errors/errors"
	"github.com/huandu/go-sqlbuilder"
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
}

func BackupDatabase(databaseName string) error {
	slog.Info("Beginning backup operation.", "log_code", "52b2d0a8", "Database-name", databaseName)

	if exists, err := util.SqliteDatabaseExists(databaseName); !exists || err != nil {
		slog.Error("Cannot add entry as database does not exist.", "log_code", "b13e3181")
		return errors.Errorf("Cannot add entry as database does not exist.")
	}

	var fileName = databaseName + ".db"

	relFs := util.NewRelativeFsManager(viper.GetString("limedbHome"))
	newDbName := fmt.Sprintf("%s-%s", fileName, strconv.FormatInt(time.Now().Unix(), 10))
	relFs.CopyFile("stores", fileName, filepath.Join("backups", databaseName), newDbName)

	db, err := util.OpenSqliteDatabase(databaseName)
	if err != nil {
		return err
	}
	defer db.Close()

	createTableStr, createTableArgs := util.CreateSqliteTable("backups", backupColumns)

	slog.Info("Creating backup table with SQL command.", "log_code", "8bc1e038", "SQL", createTableStr, "Args", createTableArgs)

	if _, err = db.Exec(createTableStr, createTableArgs...); err != nil {
		slog.Error("Failed executing create table command.", "log_code", "f7d58d42", "SQL", createTableStr)
		return err
	}

	slog.Info("Successfully created a backup table.", "log_code", "91750756", "Database-name", databaseName)

	insertStr, insertArgs := util.InsertIntoSqliteTable("backups", map[string]string{
		"date":       time.Now().Format(time.RFC3339),
		"backupName": newDbName,
	})

	slog.Info("Inserting into backup table with SQL command.", "log_code", "83d9e967", "SQL", insertStr, "Args", insertArgs)

	if _, err = db.Exec(insertStr, insertArgs...); err != nil {
		slog.Error("Failed executing insert into table command.", "log_code", "5a80e34b", "SQL", insertStr)
		return err
	}

	return RemoveOrphanBackups(databaseName)
}

func RemoveOrphanBackups(databaseName string) error {
	slog.Info("Removing backup orphans.", "log_code", "7df13463", "Database-name", databaseName)

	relFs := util.NewRelativeFsManager(viper.GetString("limedbHome"), "backups", databaseName)

	dir, err := relFs.ReadDir("")
	if err != nil {
		return err
	}

	sb := sqlbuilder.NewSelectBuilder().Distinct().Select("backupName").From("backups")
	selStr, args := sb.Build()

	db, err := util.OpenSqliteDatabase(databaseName)
	if err != nil {
		return err
	}
	defer db.Close()

    slog.Info("Querying backups in database.", "log_code", "de8e5d3d", "Database-name", databaseName, "SQL", selStr, "Args", args)

	res, err := db.Query(selStr, args...)
	if err != nil {
		slog.Warn("Could not query database backups.", "log_code", "8ddae9eb", "databaseName", databaseName)
		return err
	}
	defer res.Close()

	backupNames := make([]string, 0, 1024)

	for res.Next() {
		var val string
		res.Scan(&val)
		backupNames = append(backupNames, val)
	}

	// TODO: This is inefficient but it should be fine... How often do people backup anyway...
	for _, file := range dir {
		if !slices.Contains(backupNames, file.Name()) {
			if err := relFs.RmFile("", file.Name()); err != nil {
				return err
			}
		}
	}

	slog.Info("Successfully removed backup orphans.", "log_code", "02814e8c", "Database-name", databaseName)
	return nil
}
