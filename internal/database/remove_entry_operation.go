package database

import (
	"log/slog"

	"github.com/go-errors/errors"
	"github.com/huandu/go-sqlbuilder"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
)

func RemoveEntries(where *sqlbuilder.WhereClause) error {
	slog.Info("Beginning remove-entries operation.",
		"log_code", "406e55d9",
		"Where", where)

	selectedDb := state.ApplicationState().GetSelectedDb()

	if selectedDb == "" {
		slog.Error("Cannot add entry as no database is selected.", "log_code", "2bb83c3a")
		return errors.Errorf("Cannot add entry as no database is selected")
	}
	if exists, err := util.SqliteDatabaseExists(selectedDb); !exists || err != nil {
		slog.Error("Cannot add entry as database does not exist.", "log_code", "5e42fedb")
		return errors.Errorf("Cannot add entry as database does not exist.")
	}

	db, err := util.OpenSqliteDatabase(selectedDb)

	if err != nil {
		slog.Error("Could not open database file.", "log_code", "64a95833", "Database", selectedDb)
		return errors.Errorf("Cannot add entry as no database is selected")
	}
	defer db.Close()

    delBuilder := sqlbuilder.NewDeleteBuilder()
    delBuilder.DeleteFrom(selectedDb)
    delBuilder.AddWhereClause(where)

    sql, args := delBuilder.Build()

    slog.Info("remove-entries operation with SQL command.", "log_code", "8fc6aa55", "SQL", sql, "args", args)

	if _, err = db.Exec(sql, args...); err != nil {
		slog.Error("Failed executing remove-entries command.", "log_code", "75d6eb60", "SQL", sql, "args", args)
		return err
	}

    slog.Info("Successfully removed entries.", "log_code", "aa8cd72c", "db", selectedDb)

	return nil
}
