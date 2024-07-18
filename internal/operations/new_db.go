package operations

import (
	"database/sql"
	"io/fs"
	"log/slog"

	conf "github.com/sam-hobson/internal/config"
)

const StoresRelDir = "stores"

func NewDb(dbName string) error {
    slog.Info("Beginning new-db operation.", "Hash", "26cd37c1")

    storesDir, err := conf.GetRelDir(StoresRelDir)

    if err != nil {
        return err
    }

    fileInfo, err := fs.Stat(storesDir, dbName)

    if err != nil {
        slog.Error("Could not read/load file.", "Hash", "ff255a3d", "db-name", dbName)
        return err
    }

    db, err := sql.Open("sqlite3", fileInfo.Name())
    defer db.Close()

    slog.Info("Successfully created a new-db.", "Hash", "7bf9634b")
    return nil
}
