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

	selWhere := sqlbuilder.NewWhereClause()
	selWhere.AddWhereExpr(cond.Args, cond.Equal("rowid", rowid))

	sb := sqlbuilder.NewSelectBuilder().Select("*").From("backups").AddWhereClause(selWhere)
	selStr, selArgs := sb.Build()

	util.Log("1005f8f8").Info("restore-from-backup operation with SQL command.", "SQL", selStr, "Args", selArgs)

	res, err := db.Query(selStr, selArgs...)
	if err != nil {
		util.Log("e9290919").Error("Failed executing restore-from-backup command.", "SQL", selStr, "Args", selArgs)
		return err
	}
	defer res.Close()

	var date string
	var backupName string
	var comment string

	if !res.Next() {
		util.Log("cf85bbc9").Error("Failed restore-from-backup, no matching backup.")
		return errors.Errorf("Failed restore-from-backup, no matching backup")
	}

	res.Scan(&date, &backupName, &comment)

	if res.Next() {
		util.Log("16fc100d").Error("Failed restore-from-backup, too many backup field records.", "Backup name retrieved", backupName)
		return errors.Errorf("Failed restore-from-backup, too many backup field records")
	}

    fileName := databaseName + ".db"
	relFs := util.NewRelativeFsManager(state.ApplicationState().GetLimedbHome())
	if err := relFs.CopyFile(filepath.Join("backups", databaseName, backupName), filepath.Join("stores", fileName)); err != nil {
		return err
	}

	// Open the newly restored database
	db, err = dbutil.OpenSqliteDatabaseIfExists(databaseName)
	if err != nil {
		return err
	}
	defer db.Close()

    util.Log("e2ab58c3").Info("Successfully restored backup.", "Database name", databaseName, "Row id", rowid)

	return nil
}
