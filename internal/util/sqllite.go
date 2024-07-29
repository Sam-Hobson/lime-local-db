package util

import (
	"database/sql"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func SqliteDatabaseExists(dbName string) (bool, error) {
	fileName := dbName

	if !strings.HasSuffix(fileName, ".db") {
		fileName += ".db"
	}

	slog.Info("Checking if database exists.", "log_code", "e75f8412", "dbName", dbName)

	return NewRelativeFsManager(viper.GetString("limedbHome")).FileExists("stores", fileName)
}

func OpenSqliteDatabase(dbName string) (*sql.DB, error) {
	fileName := dbName

	if !strings.HasSuffix(fileName, ".db") {
		fileName += ".db"
	}
	dbPath := filepath.Join(viper.GetString("limedbHome"), "stores", fileName)

	slog.Info("Opening database file.", "log_code", "34503562", "db_path", dbPath)

	return sql.Open("sqlite3", dbPath)
}

func AllExistingDatabaseNames() []string {

    return nil
}
