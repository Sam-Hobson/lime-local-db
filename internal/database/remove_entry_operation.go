package database

import (
	"log/slog"

	"github.com/go-errors/errors"
	"github.com/huandu/go-sqlbuilder"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
)

func RemoveEntries(where *sqlbuilder.WhereClause) (int64, error) {
	slog.Info("Beginning remove-entries operation.",
		"Log code", "406e55d9",
		"Where", where)

	selectedDb := state.ApplicationState().GetSelectedDb()

	if selectedDb == "" {
		slog.Error("Cannot add entry as no database is selected.", "Log code", "2bb83c3a")
		return -1, errors.Errorf("Cannot add entry as no database is selected")
	}
	if exists, err := util.SqliteDatabaseExists(selectedDb); !exists || err != nil {
		slog.Error("Cannot add entry as database does not exist.", "Log code", "5e42fedb")
		return -1, errors.Errorf("Cannot add entry as database does not exist.")
	}

	db, err := util.OpenSqliteDatabase(selectedDb)

	if err != nil {
		slog.Error("Could not open database file.", "Log code", "64a95833", "Database", selectedDb)
		return -1, errors.Errorf("Cannot add entry as no database is selected")
	}
	defer db.Close()

	delBuilder := sqlbuilder.NewDeleteBuilder()
	delBuilder.DeleteFrom(selectedDb)
	delBuilder.AddWhereClause(where)

	sql, args := delBuilder.Build()
	slog.Info("remove-entries operation with SQL command.", "Log code", "8fc6aa55", "SQL", sql, "Args", args)
	res, err := db.Exec(sql, args...)

	if err != nil {
		slog.Error("Failed executing remove-entries command.", "Log code", "75d6eb60", "SQL", sql, "Args", args)
		return -1, err
	}

	rowsAffected := util.PanicIfErr(res.RowsAffected())

	slog.Info("Successfully removed entries.", "Log code", "aa8cd72c", "Rows affected", rowsAffected, "Db", selectedDb)

	return rowsAffected, nil
}
