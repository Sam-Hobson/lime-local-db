package operations

import (
	"log/slog"
	"strings"

	cp "github.com/bigkevmcd/go-configparser"
	"github.com/sam-hobson/internal/config"
)

const SoftDeleteFileDir = "deleted"

func RmDb(dbName string, conf *cp.ConfigParser) error {
	slog.Info("Beginning rm-db operation.", "log_code", "49aaf185", "db-name", dbName)

    if !strings.HasSuffix(dbName, ".db") {
        dbName += ".db"
    }

    shouldRm, err := conf.GetBool("STORE", "permanently_delete_on_rm")

    if err != nil {
        slog.Warn("Could not find \"[STORE]\"\"permanently_delete_on_rm\" in config file. Defaulting to false.", "log_code", "77c34305")
        shouldRm = false
    }

    if shouldRm {
        if err := config.RmFile(StoresRelDir, dbName); err != nil {
            return err
        }
    } else {
        if err := config.MoveFile(StoresRelDir, dbName, SoftDeleteFileDir, dbName); err != nil {
            return err
        }
    }

	slog.Info("Successfully completed rm-db operation.", "log_code", "d73a061e", "db-name", dbName)
    return nil
}
