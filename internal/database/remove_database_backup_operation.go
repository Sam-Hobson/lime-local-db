package database

import (
	"github.com/go-errors/errors"
	"github.com/huandu/go-sqlbuilder"
	dbutil "github.com/sam-hobson/internal/database/util"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/viper"
)

func RemoveDatabaseBackup(databaseName string, rowid int) error {
	util.Log("01fcf023").Info("Beginning remove backup operation.", "Database name", databaseName, "Row id", rowid)

	persistentDatabaseName := dbutil.PersistentDatabaseName(databaseName)

	db, err := dbutil.OpenSqliteDatabaseIfExists(persistentDatabaseName)
	if err != nil {
		return err
	}
	defer db.Close()

	cond := sqlbuilder.NewCond()

	selWhere := sqlbuilder.NewWhereClause()
	selWhere.AddWhereExpr(cond.Args, cond.Equal("rowid", rowid))

	sb := sqlbuilder.NewSelectBuilder().Select("backupName").From("backups").AddWhereClause(selWhere)
	selStr, selArgs := sb.Build()

	util.Log("b5636cf9").Info("remove-database-backup operation with SQL command.", "SQL", selStr, "Args", selArgs)

	res, err := db.Query(selStr, selArgs...)
	if err != nil {
		util.Log("76116735").Error("Failed executing remove-database-backup command.", "SQL", selStr, "Args", selArgs)
		return err
	}
	defer res.Close()

	var backupName string

	if !res.Next() {
		util.Log("e8f55bfc").Error("Failed remove-database-backup, no matching backup.")
		return errors.Errorf("Failed remove-database-backup, no matching backup")
	}

	res.Scan(&backupName)

	if res.Next() {
		util.Log("734c3115").Error("Failed remove-database-backup, too many backup field records.", "Backup name retrieved", backupName)
		return errors.Errorf("Failed remove-database-backup, too many backup field records")
	}

	delWhere := sqlbuilder.NewWhereClause()
	delWhere.AddWhereExpr(cond.Args, cond.Equal("rowid", rowid))

	delBuilder := sqlbuilder.NewDeleteBuilder().DeleteFrom("backups").AddWhereClause(delWhere)
	delStr, delArgs := delBuilder.Build()

	util.Log("b94d03eb").Info("remove-database-backup operation with SQL command.", "SQL", delStr, "Args", delArgs)

	if _, err := db.Exec(delStr, delArgs...); err != nil {
		util.Log("08a4a0df").Error("Failed remove-database-backup, could not delete backup.", "SQL", delStr, "Args", delArgs)
		return errors.Errorf("Failed remove-database-backup, could not delete backup")
	}

	relFs := util.NewRelativeFsManager(viper.GetString("limedb_home"), "backups")

	if err := relFs.RmFile(databaseName, backupName); err != nil {
		return err
	}

	util.Log("f1fe1bdb").Info("Successfully removed backup.", "Persistent database name", persistentDatabaseName, "Row id", rowid, "Backup name", backupName)
	return nil
}
