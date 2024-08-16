package util

import (
	"database/sql"

	"github.com/go-errors/errors"
	"github.com/huandu/go-sqlbuilder"
	"github.com/sam-hobson/internal/util"
)

type Trigger struct {
	Name        string
	Sql         string
	DateCreated string
	TriggerType string
	Comment     string
}

func DefinedTriggers(databaseName string) ([]*Trigger, error) {
	util.Log("89518643").Info("Getting defined triggers for table.", "Database name", databaseName)

	db, err := OpenSqliteDatabaseIfExists(databaseName)
	if err != nil {
		return nil, err
	}

	cond := sqlbuilder.NewCond()
	selStr, args := EntriesInTableWhereSql("sqlite_master", []string{"name", "sql"}, cond.Args, cond.Equal("type", "trigger"))

	util.Log("c5de02bb").Info("Querying defined triggers.", "Database name", databaseName, "SQL", selStr, "Args", args)

	res, err := db.Query(selStr, args...)
	if err != nil {
		util.Log("5cd94c57").Warn("Could not query defined triggers.", "Database name", databaseName)
		return nil, err
	}
	defer res.Close()

	triggers := make([]*Trigger, 0)

	for res.Next() {
		trigger := Trigger{}
		res.Scan(&trigger.Name, &trigger.Sql)
		triggers = append(triggers, &trigger)
	}

	util.Log("35fa533b").Info("Defined triggers found.", "Database name", databaseName, "Triggers", triggers)

	return triggers, nil
}

func TriggerNameExists(db *sql.DB, triggerName string) (bool, error) {
	util.Log("be8c5e04").Info("Check if trigger exists with name", "Trigger name", triggerName)
	cond := sqlbuilder.NewCond()
	sql, args := EntriesInTableWhereSql("triggers", []string{"*"}, cond.Args, cond.Equal("name", triggerName))

	util.Log("fa8d9b89").Info("Querying triggers table with SQL.", "SQL", sql, "Args", args)

	if res, err := db.Exec(sql, args...); err != nil {
		util.Log("13bd616e").Warn("Failed to query triggers table.")
		return true, err
	} else if util.PanicIfErr(res.RowsAffected()) != 0 {
		return true, nil
	}

	return false, nil
}
