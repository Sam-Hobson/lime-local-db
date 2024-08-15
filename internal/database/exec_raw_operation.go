package database

import (
	"database/sql"

	dbutil "github.com/sam-hobson/internal/database/util"
	"github.com/sam-hobson/internal/util"
)

func ExecRawSql(databaseName, sql string) (sql.Result, error) {
	util.Log("81a1e54f").Info("Beginning execute raw sql operation.", "Database name", databaseName, "SQL", sql)

	db, err := dbutil.OpenSqliteDatabaseIfExists(databaseName)
	if err != nil {
		return nil, err
	}
	defer db.Close()

    util.Log("5698ba36").Info("Executing SQL on database.", "Database name", databaseName, "SQL", sql)

	res, err := db.Exec(sql)
	if err != nil {
		util.Log("b3b635d0").Error("Could not execute raw sql contents on database.", "Database name", databaseName)
		return nil, err
	}

	return res, nil
}
