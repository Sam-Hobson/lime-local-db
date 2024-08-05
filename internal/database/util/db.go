package util

import (
	"fmt"

	"github.com/sam-hobson/internal/util"
	"github.com/spf13/viper"
)

// TODO: This should be refactored into a struct
func AllExistingDatabaseNames() ([]string, error) {
	relFs := util.NewRelativeFsManager(viper.GetString("limedb_home"))
	files, err := relFs.ReadDir("stores")

	if err != nil {
		return nil, err
	}

	dbNames := make([]string, len(files))

	for i, file := range files {
		dbNames[i] = file.Name()
	}

	return dbNames, nil
}

func PersistentDatabaseName(databaseName string) string {
    return fmt.Sprintf(".%s_persistent", databaseName)
}
