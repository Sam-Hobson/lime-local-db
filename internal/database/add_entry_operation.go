package database

import (
	"github.com/sam-hobson/internal/util"
	dbutil "github.com/sam-hobson/internal/database/util"
)

func AddEntry(databaseName string, entries map[string]string) error {
	util.Log("f3f1b8df").Info("Beginning add-entry operation.", "Entries", entries)

    db, err := dbutil.OpenSqliteDatabaseIfExists(databaseName)
	if err != nil {
        return err
	}
	defer db.Close()

	insertStr, args := dbutil.InsertIntoTableSql(databaseName, entries)

	util.Log("01809774").Info("Inserting with SQL Command.", "SQL", insertStr, "Args", args)

	if _, err = db.Exec(insertStr, args...); err != nil {
		util.Log("0981c049").Error("Failed executing insert table command.", "SQL", insertStr)
		return err
	}

	util.Log("3e11ab9a").Info("Successfully inserted into database.", "Database name", databaseName)

	return nil
}
