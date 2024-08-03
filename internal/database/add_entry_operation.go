package database

import (
	"log/slog"

	"github.com/go-errors/errors"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
)

func AddEntry(entries map[string]string) error {
	slog.Info("Beginning add-entry operation.",
		"log_code", "f3f1b8df",
		"Entries", entries)

	selectedDb := state.ApplicationState().GetSelectedDb()

	if selectedDb == "" {
		slog.Error("Cannot add entry as no database is selected.", "log_code", "c40be9f9")
		return errors.Errorf("Cannot add entry as no database is selected")
	}
	if exists, err := util.SqliteDatabaseExists(selectedDb); !exists || err != nil {
		slog.Error("Cannot add entry as database does not exist.", "log_code", "765a9254")
		return errors.Errorf("Cannot add entry as database does not exist.")
	}

	db, err := util.OpenSqliteDatabase(selectedDb)
	if err != nil {
		return err
	}
	defer db.Close()

	insertStr, args := util.InsertIntoSqliteTable(selectedDb, entries)

	slog.Info("Inserting with SQL Command.",
		"log_code", "01809774",
		"SQL", insertStr,
		"args", args)

	if _, err = db.Exec(insertStr, args...); err != nil {
		slog.Error("Failed executing insert table command.", "log_code", "0981c049", "SQL", insertStr)
		return err
	}

	slog.Info("Successfully inserted into database.", "log_code", "3e11ab9a", "Selected_db", selectedDb)

	return nil
}
