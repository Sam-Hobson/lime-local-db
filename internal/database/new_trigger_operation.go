package database

import (
	dbutil "github.com/sam-hobson/internal/database/util"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
)

func NewTriggerFromFile(databaseName string, fileName string) error {
	util.Log("d8e9dc08").Info("Beginning new trigger operation.", "Database name", databaseName, "File name", fileName)

	relFs := util.NewRelativeFsManager(state.ApplicationState().GetLimedbHome(), "triggers", databaseName)
	fileContents, err := relFs.ReadFileIntoMemry("", fileName)

	if err != nil {
		util.Log("148e9567").Error("Could not source trigger file.", "File name", fileName)
		return err
	}

	db, err := dbutil.OpenSqliteDatabaseIfExists(databaseName)
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err := db.Exec(fileContents); err != nil {
		util.Log("4fd6ac3a").Error("Could not execute trigger contents on database.", "Database name", databaseName, "File name", fileName)
		return err
	}

    util.Log("49d0ff3d").Info("Successfully created trigger from file.", "Database name", databaseName, "File name", fileName)

	return nil
}
