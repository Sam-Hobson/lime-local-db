package operations

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/go-errors/errors"
	_ "github.com/mattn/go-sqlite3"

	conf "github.com/sam-hobson/internal/config"
)

const StoresRelDir = "stores"

type NewdbData struct {
	Dbname   string
	ColNames []string
}

func (f *NewdbData) String() string {
	return fmt.Sprintf("%+v", *f)
}

func NewDb(cmd *NewdbData) error {
	slog.Info("Beginning new-db operation.", "log_code", "26cd37c1", "cmd", cmd)

	dbName := cmd.Dbname

	if !strings.HasSuffix(dbName, ".db") {
		dbName += ".db"
	}

	exists, err := conf.FileExists(StoresRelDir, dbName)

	if err != nil {
		return err
	}

	if exists {
		slog.Error("Cannot create a new database as it already exists.", "log_code", "6c95edf6", "path", conf.FullPath(StoresRelDir, dbName))
		return errors.Errorf("Cannot create a new database as it already exists.")
	}

	err = conf.CreateFile(StoresRelDir, dbName)

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

    // TODO: Insert into db

	slog.Info("Successfully created a new-db.", "log_code", "7bf9634b", "Db_path", dbPath)
	return nil
}
