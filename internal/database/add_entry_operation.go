package database

import (
	"github.com/go-errors/errors"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
	dbutil "github.com/sam-hobson/internal/database/util"
)

func AddEntry(entries map[string]string) error {
	util.Log("f3f1b8df").Info("Beginning add-entry operation.", "Entries", entries)
	selectedDb := state.ApplicationState().GetSelectedDb()

	if selectedDb == "" {
		util.Log("c40be9f9").Error("Cannot add entry as no database is selected.")
		return errors.Errorf("Cannot add entry as no database is selected")
	}

    db, err := dbutil.OpenSqliteDatabaseIfExists(selectedDb)
	if err != nil {
        return err
	}
	defer db.Close()

	insertStr, args := dbutil.InsertIntoTableSql(selectedDb, entries)

	util.Log("01809774").Info("Inserting with SQL Command.", "SQL", insertStr, "Args", args)

	if _, err = db.Exec(insertStr, args...); err != nil {
		util.Log("0981c049").Error("Failed executing insert table command.", "SQL", insertStr)
		return err
	}

	util.Log("3e11ab9a").Info("Successfully inserted into database.", "Selected database", selectedDb)

	return nil
}
