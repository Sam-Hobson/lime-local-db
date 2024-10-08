package database

import (
	"fmt"
	"path/filepath"
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

		relFs.MoveFile(filepath.Join("stores", fileName), filepath.Join("deleted", newDbName))
		relFs.MoveFile(filepath.Join("stores", persistentFileName), filepath.Join("deleted", newPersistentDbName))
	} else {
		relFs.RmFile("stores", fileName)
		relFs.RmFile("stores", persistentFileName)
	}

	util.Log("d73a061e").Info("db rm operation executed successfully.", "Database name", fileName, "Soft delete", softDelete)
	return nil

}
