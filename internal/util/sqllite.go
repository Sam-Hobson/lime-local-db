package util

import (
	"database/sql"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func OpenSqliteDatabase(dbName string) (*sql.DB, error) {
	fileName := dbName

	if !strings.HasSuffix(fileName, ".db") {
		fileName += ".db"
	}

	dbPath := filepath.Join(viper.GetString("limedbHome"), "stores", fileName)

	slog.Info("Opening database file.", "log_code", "34503562", "db_path", dbPath)

	return sql.Open("sqlite3", dbPath)
}
