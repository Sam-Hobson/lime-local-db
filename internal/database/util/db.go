package util

import (
	"fmt"

	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
)

// TODO: This should be refactored into a struct
func AllExistingDatabaseNames() ([]string, error) {
	relFs := util.NewRelativeFsManager(state.ApplicationState().GetLimedbHome())
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
