package util

import (
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
