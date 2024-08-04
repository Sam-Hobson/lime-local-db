package database

import (
	"fmt"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/sam-hobson/internal/types"
	"github.com/sam-hobson/internal/util"
	dbutil "github.com/sam-hobson/internal/database/util"
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

func BackupDatabase(databaseName, comment string) error {
	util.Log("52b2d0a8").Info("Beginning backup operation.", "Database name", databaseName)

	db, err := dbutil.OpenSqliteDatabaseIfExists(databaseName)
	if err != nil {
		return err
	}
	defer db.Close()

	fileName := databaseName + ".db"
	relFs := util.NewRelativeFsManager(viper.GetString("limedb_home"))
	newDbName := fmt.Sprintf("%s-%s", fileName, strconv.FormatInt(time.Now().Unix(), 10))

	relFs.CopyFile("stores", fileName, filepath.Join("backups", databaseName), newDbName)

	createTableStr, createTableArgs := dbutil.CreateSqliteTable("backups", backupColumns)

	util.Log("8bc1e038").Info("Creating backup table with SQL command.", "SQL", createTableStr, "Args", createTableArgs)

	if _, err = db.Exec(createTableStr, createTableArgs...); err != nil {
		util.Log("f7d58d42").Error("Failed executing create table command.", "SQL", createTableStr)
		return err
	}

	util.Log("91750756").Info("Successfully created a backup table.", "Database name", databaseName)

	insertStr, insertArgs := dbutil.InsertIntoSqliteTable("backups", map[string]string{
		"date":       time.Now().Format(time.RFC3339),
		"backupName": newDbName,
		"comment":    comment,
	})

	util.Log("83d9e967").Info("Inserting into backup table with SQL command.", "SQL", insertStr, "Args", insertArgs)

	if _, err = db.Exec(insertStr, insertArgs...); err != nil {
		util.Log("5a80e34b").Error("Failed executing insert into table command.", "SQL", insertStr)
		return err
	}

	if viper.GetBool("remove_orphan_backups") {
		return RemoveOrphanBackups(databaseName)
	}

	return nil
}

func RemoveOrphanBackups(databaseName string) error {
	util.Log("7df13463").Info("Removing backup orphans.", "Database name", databaseName)

	relFs := util.NewRelativeFsManager(viper.GetString("limedb_home"), "backups", databaseName)

	dir, err := relFs.ReadDir("")
	if err != nil {
		return err
	}

	sb := sqlbuilder.NewSelectBuilder().Distinct().Select("backupName").From("backups")
	selStr, args := sb.Build()

	db, err := dbutil.OpenSqliteDatabaseIfExists(databaseName)
	if err != nil {
		return err
	}
	defer db.Close()

	util.Log("de8e5d3d").Info("Querying backups in database.", "Database name", databaseName, "SQL", selStr, "Args", args)

	res, err := db.Query(selStr, args...)
	if err != nil {
		util.Log("8ddae9eb").Warn("Could not query database backups.", "Database name", databaseName)
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

	util.Log("02814e8c").Info("Successfully removed backup orphans.", "Database name", databaseName)
	return nil
}
