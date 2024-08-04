package database

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/sam-hobson/internal/util"
	"github.com/spf13/viper"
)

func RemoveDatabase(databaseName string) error {
	slog.Info("Beginning rm-db operation.", "Log code", "49aaf185", "Database name", databaseName)

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

	slog.Info("Successfully removed database.", "Log code", "d73a061e", "Db name", fileName, "Soft delete", softDelete)
	return nil

}
