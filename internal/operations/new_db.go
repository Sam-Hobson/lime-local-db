package operations

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	conf "github.com/sam-hobson/internal/config"
)

const StoresRelDir = "stores"

type NewdbCmd struct {
	Dbname   string
	ColNames []string
}

func (f *NewdbCmd) String() string {
	return fmt.Sprintf("%+v", *f)
}

func NewDb(cmd *NewdbCmd) error {
	slog.Info("Beginning new-db operation.", "log_code", "26cd37c1", "cmd", cmd)

    dbName := cmd.Dbname

	if !strings.HasSuffix(dbName, ".db") {
		dbName += ".db"
	}

	// Create the dir if it doesn't exist
	err := conf.CreateFile(StoresRelDir, dbName)

	if err != nil {
		return err
	}

	dbPath := conf.FullPath(StoresRelDir, dbName)
	db, err := sql.Open("sqlite3", dbPath)
	defer db.Close()

	if err != nil {
		slog.Error("Could not open sqlite database.", "log_code", "9494fc60", "Db_path", dbPath)
		return err
	}

	// TODO: Actually create the database.

	slog.Info("Successfully created a new-db.", "log_code", "7bf9634b", "Db_path", dbPath)
	return nil
}
