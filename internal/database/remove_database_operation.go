package database

import (
	"fmt"
	"strconv"
	"time"

	"github.com/sam-hobson/internal/util"
	"github.com/spf13/viper"
)

func RemoveDatabase(databaseName string) error {
	util.Log("49aaf185").Info("Beginning rm-db operation.", "Database name", databaseName)

	fileName := databaseName + ".db"

	relFs := util.NewRelativeFsManager(viper.GetString("limedb_home"))
	softDelete := viper.GetBool("soft_deletion")

	if softDelete {
		newDbName := fmt.Sprintf("%s-%s", fileName, strconv.FormatInt(time.Now().Unix(), 10))
		if err := relFs.MoveFile("stores", fileName, "deleted", newDbName); err != nil {
			return err
		}
	} else {
		if err := relFs.RmFile("stores", fileName); err != nil {
			return err
		}
	}

	util.Log("d73a061e").Info("Successfully removed database.", "Db name", fileName, "Soft delete", softDelete)
	return nil

}
