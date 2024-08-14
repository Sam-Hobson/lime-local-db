package database

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	dbutil "github.com/sam-hobson/internal/database/util"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
)

func BackupDatabase(databaseName, comment string) (int64, error) {
	util.Log("52b2d0a8").Info("Beginning backup operation.", "Database name", databaseName)

	persistentDatabaseName := dbutil.PersistentDatabaseName(databaseName)

	db, err := dbutil.OpenSqliteDatabaseIfExists(persistentDatabaseName)
	if err != nil {
		return -1, err
	}
	defer db.Close()

	fileName := databaseName + ".db"
	relFs := util.NewRelativeFsManager(state.ApplicationState().GetLimedbHome())
	newDbName := fmt.Sprintf("%s-%s", fileName, strconv.FormatInt(time.Now().Unix(), 10))

	relFs.CopyFile(filepath.Join("stores", fileName), filepath.Join("backups", databaseName, newDbName))

	insertStr, insertArgs := dbutil.InsertIntoTableSql("backups", map[string]string{
		"date":       time.Now().Format(time.RFC3339),
		"backupName": newDbName,
		"comment":    comment,
	})

	util.Log("83d9e967").Info("Inserting into backup table with SQL command.", "SQL", insertStr, "Args", insertArgs)

	res, err := db.Exec(insertStr, insertArgs...)
	if err != nil {
		util.Log("5a80e34b").Error("Failed executing insert into table command.", "SQL", insertStr)
		return -1, err
	}

	return util.PanicIfErr(res.LastInsertId()), nil
}
