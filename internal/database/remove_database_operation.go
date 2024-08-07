package database

import (
	"fmt"
	"strconv"
	"time"

	dbutil "github.com/sam-hobson/internal/database/util"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/viper"
)

func RemoveDatabase(databaseName string) error {
	util.Log("49aaf185").Info("Beginning db rm operation.", "Database name", databaseName)

	fileName := databaseName + ".db"
	persistentFileName := dbutil.PersistentDatabaseName(databaseName) + ".db"

	relFs := util.NewRelativeFsManager(state.ApplicationState().GetLimedbHome())
	softDelete := viper.GetBool("soft_deletion")

	if softDelete {
		currentTimestamp := strconv.FormatInt(time.Now().Unix(), 10)
		newDbName := fmt.Sprintf("%s-%s", fileName, currentTimestamp)
		newPersistentDbName := fmt.Sprintf("%s-%s", persistentFileName, currentTimestamp)

		if err := relFs.MoveFile("stores", fileName, "deleted", newDbName); err != nil {
			util.Log("53913fa0").Warn("Could not soft delete database.", "Database name", fileName)
		}
		if err := relFs.MoveFile("stores", persistentFileName, "deleted", newPersistentDbName); err != nil {
			util.Log("3269a471").Warn("Could not soft delete persistent database.", "Persistent database name", persistentFileName)
		}
	} else {
		if err := relFs.RmFile("stores", fileName); err != nil {
			util.Log("5d32fcdd").Warn("Could not delete database.", "Database name", fileName)
		}
		if err := relFs.RmFile("stores", persistentFileName); err != nil {
			util.Log("2560d564").Warn("Could not delete persistent database.", "Persistent database name", persistentFileName)
		}
	}

	util.Log("d73a061e").Info("db rm operation executed successfully.", "Db name", fileName, "Soft delete", softDelete)
	return nil

}
