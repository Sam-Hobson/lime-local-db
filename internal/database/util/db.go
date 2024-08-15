package util

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
)

// TODO: This should be refactored into a struct
func AllExistingDatabaseNames() ([]string, error) {
	relFs := util.NewRelativeFsManager(state.ApplicationState().GetLimedbHome())
	files, err := relFs.ReadDir("stores")

	if err != nil {
		return nil, err
	}

	dbNames := make([]string, 0)

	for _, file := range files {
		name := file.Name()

		if !strings.HasPrefix(name, ".") {
			if strings.HasSuffix(name, ".db") {
				dbNames = append(dbNames, name[:len(name)-3])
			} else {
				dbNames = append(dbNames, name)
			}
		}
	}

	return dbNames, nil
}

func PersistentDatabaseName(databaseName string) string {
	return fmt.Sprintf(".%s_persistent", databaseName)
}

func OpenSqliteDatabaseIfExists(databaseName string) (*sql.DB, error) {
	exists, err := SqliteDatabaseExists(databaseName)

	if err != nil {
		util.Log("94f1ece2").Warn("Cannot open database as it does not exist.")
		return nil, err
	}
	if !exists {
		util.Log("cbd713ec").Warn("Cannot open database as it does not exist.")
		return nil, errors.Errorf("Cannot open database as it does not exist")
	}

	return OpenSqliteDatabase(databaseName)
}

func SqliteDatabaseExists(databaseName string) (bool, error) {
	fileName := databaseName + ".db"
	relFs := util.NewRelativeFsManager(state.ApplicationState().GetLimedbHome(), "stores")
	util.Log("e75f8412").Info("Checking if database exists.", "Database name", databaseName)
	return relFs.FileExists(fileName)
}

func OpenSqliteDatabase(databaseName string) (*sql.DB, error) {
	fileName := databaseName + ".db"
	dbPath := filepath.Join(state.ApplicationState().GetLimedbHome(), "stores", fileName)
	util.Log("34503562").Info("Opening database file.", "Database path", dbPath)
	return sql.Open("sqlite3", dbPath)
}

func RemoveSqliteDatabase(databaseName string) error {
	fileName := databaseName + ".db"
	relFs := util.NewRelativeFsManager(state.ApplicationState().GetLimedbHome(), "stores")
	util.Log("cb25f7f8").Info("Removing database file.", "Database name", databaseName)
	return relFs.RmFile(fileName)
}
