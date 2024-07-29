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

// TODO: This should be refactored into a struct
func AllExistingDatabaseNames() ([]string, error) {
	files, err := NewRelativeFsManager(viper.GetString("limedbHome")).ReadDir("stores")

	if err != nil {
		return nil, err
	}

	dbNames := make([]string, len(files))

	for i, file := range files {
		dbNames[i] = file.Name()
	}

	return dbNames, nil
}
