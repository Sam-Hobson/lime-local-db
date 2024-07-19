package operations

import (
	"log/slog"
	"strings"

	"github.com/sam-hobson/internal/config"
)

const SoftDeleteFileDir = "deleted"

func RmDb(dbName string) error {
	slog.Info("Beginning rm-db operation.", "log_code", "49aaf185", "db-name", dbName)

	conf, err := config.GetConfig()

	if err != nil {
		return err
	}

	if !strings.HasSuffix(dbName, ".db") {
		dbName += ".db"
	}

	softDelete, err := conf.GetBool("STORE", "soft_delete_on_rm")

	if err != nil {
		slog.Warn("Could not find \"[STORE]\"\"soft_delete_on_rm\" in config file. Defaulting to false.", "log_code", "77c34305")
		softDelete = true
	}

	if softDelete {
		if err := config.MoveFile(StoresRelDir, dbName, SoftDeleteFileDir, dbName); err != nil {
			return err
		}
	} else {
		if err := config.RmFile(StoresRelDir, dbName); err != nil {
			return err
		}
	}

	slog.Info("Successfully completed rm-db operation.", "log_code", "d73a061e", "db-name", dbName)
	return nil
}
