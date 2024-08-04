package database

import (
	"github.com/go-errors/errors"
	"github.com/huandu/go-sqlbuilder"
	dbutil "github.com/sam-hobson/internal/database/util"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
)

func RemoveEntries(where *sqlbuilder.WhereClause) (int64, error) {
	util.Log("406e55d9").Info("Beginning remove-entries operation.", "Where", where)
	selectedDb := state.ApplicationState().GetSelectedDb()

	if selectedDb == "" {
		util.Log("2bb83c3a").Error("Cannot add entry as no database is selected.")
		return -1, errors.Errorf("Cannot add entry as no database is selected")
	}

	db, err := dbutil.OpenSqliteDatabaseIfExists(selectedDb)
	if err != nil {
		util.Log("64a95833").Error("Could not open database file.", "Database", selectedDb)
		return -1, errors.Errorf("Cannot add entry as no database is selected")
	}
	defer db.Close()

	delBuilder := sqlbuilder.NewDeleteBuilder().DeleteFrom(selectedDb).AddWhereClause(where)
	sql, args := delBuilder.Build()
	util.Log("8fc6aa55").Info("remove-entries operation with SQL command.", "SQL", sql, "Args", args)

	res, err := db.Exec(sql, args...)
	if err != nil {
		util.Log("75d6eb60").Error("Failed executing remove-entries command.", "SQL", sql, "Args", args)
		return -1, err
	}

	rowsAffected := util.PanicIfErr(res.RowsAffected())

	util.Log("aa8cd72c").Info("Successfully removed entries.", "Rows affected", rowsAffected, "Db", selectedDb)

	return rowsAffected, nil
}
