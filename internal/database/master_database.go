package database

import (
	"database/sql"
	"path/filepath"

	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/types"
	"github.com/sam-hobson/internal/util"
)

const MasterDatabaseName = ".master"
const MasterDatabaseFileName = ".master.db"

// func MasterDatabaseExists() (bool, error) {
// 	relFs := util.NewRelativeFsManager(state.ApplicationState().GetLimedbHome())
//
// 	if exists, err := relFs.FileExists("", MasterDatabaseFileName); err != nil {
// 		return false, err
// 	} else {
// 		return exists, nil
// 	}
// }
//
// func CreateMasterDatabase() error {
// 	if exists, err := MasterDatabaseExists(); err != nil {
// 		return err
// 	} else if exists {
// 		util.Log("0618d59f").Warn("Attempted to create master database when it already exists.")
// 		return errors.Errorf("Attempted to create master database when it already exists")
// 	}
//
// 	relFs := util.NewRelativeFsManager(state.ApplicationState().GetLimedbHome())
//
// 	if err := relFs.CreateFile("", MasterDatabaseName); err != nil {
// 		return err
// 	} else {
// 		util.Log("9e4e0299").Info("Successfully created master database.")
// 		return nil
// 	}
// }

func AddNewDatabaseRecord(databaseName string, columns []*types.Column) error {
    util.Log("ba44622e").Info("Adding new database record to master database.", "Database name", databaseName, "Columns", columns)

    db, err := OpenMasterDatabase()
    if err != nil {
        return err
    }
    defer db.Close()

    return nil
}

func OpenMasterDatabase() (*sql.DB, error) {
    if db, err := sql.Open("sqlite3", MasterDatabaseFilePath()); err != nil {
        util.Log("facf46fb").Warn("Could not open master database.")
        return nil, err
    } else {
        return db, nil
    }
}

func MasterDatabaseFilePath() string {
	return filepath.Join(state.ApplicationState().GetLimedbHome(), MasterDatabaseFileName)
}
