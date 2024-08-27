package database

import (
	"path/filepath"

	"github.com/go-errors/errors"
	"github.com/huandu/go-sqlbuilder"
	dbutil "github.com/sam-hobson/internal/database/util"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
)

func RestoreFromBackup(databaseName string, rowid int) error {
	util.Log("f68180a2").Info("Beginning restore from backup operation.", "Database name", databaseName, "Row id", rowid)

	persistentDatabaseName := dbutil.PersistentDatabaseName(databaseName)

	db, err := dbutil.OpenSqliteDatabaseIfExists(persistentDatabaseName)
	if err != nil {
		return err
	}
	defer db.Close()

	cond := sqlbuilder.NewCond()
	selStr, selArgs := dbutil.EntriesInTableWhereSql("backups", []string{"name"}, cond.Args, cond.Equal("rowid", rowid))

	util.Log("1005f8f8").Info("restore-from-backup operation with SQL command.", "SQL", selStr, "Args", selArgs)

	res, err := db.Query(selStr, selArgs...)
	if err != nil {
		util.Log("e9290919").Error("Failed executing restore-from-backup command.", "SQL", selStr, "Args", selArgs)
		return err
	}
	defer res.Close()

	names, err := dbutil.RowsIntoSlice[string](res)
	if err != nil {
		util.Log("a122ca2e").Error("Could not read backup into slice.", "Database name", databaseName, "Backup row id", rowid)
		return err
	}

	if len(names) == 0 {
		util.Log("cf85bbc9").Error("Failed restore-from-backup, no matching backup.")
		return errors.Errorf("Failed restore-from-backup, no matching backup")
	} else if len(names) > 1 {
		util.Log("16fc100d").Error("Failed restore-from-backup, too many backup field records.", "Backup names", names)
		return errors.Errorf("Failed restore-from-backup, too many backup field records")
	}

	backupName := names[0]

	fileName := databaseName + ".db"
	relFs := util.NewRelativeFsManager(state.ApplicationState().GetLimedbHome())
	if err := relFs.CopyFile(filepath.Join("backups", databaseName, backupName), filepath.Join("stores", fileName)); err != nil {
		return err
	}

	util.Log("e2ab58c3").Info("Successfully restored backup.", "Database name", databaseName, "Row id", rowid)

	return nil
}
