package database

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/go-errors/errors"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/viper"
)

func BackupDatabase(databaseName string) error {
    slog.Info("Beginning backup operation.", "log_code", "52b2d0a8", "Database-name", databaseName)

    if exists, err := util.SqliteDatabaseExists(databaseName); !exists || err != nil {
        slog.Error("Cannot add entry as database does not exist.", "log_code", "b13e3181")
        return errors.Errorf("Cannot add entry as database does not exist.")
    }

    var fileName = databaseName
    if !strings.HasSuffix(databaseName, ".db") {
        fileName += ".db"
    }

    relFs := util.NewRelativeFsManager(viper.GetString("limedbHome"))

    newDbName := fmt.Sprintf("%s-%s", fileName, strconv.FormatInt(time.Now().Unix(), 10))
    relFs.CopyFile("stores", fileName, "backups", newDbName)

    db, err := util.OpenSqliteDatabase(databaseName)
    if err != nil {
        return err
    }
    defer db.Close()


    return nil
}
