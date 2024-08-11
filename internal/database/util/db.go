package util

import (
	"fmt"
	"strings"

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

	dbNames := make([]string, 0)

	for _, file := range files {
		name := file.Name()

		if !strings.HasPrefix(name, ".") {
			if strings.HasSuffix(name, ".db") {
				dbNames = append(dbNames, name[:len(name)-3])
			} else {
				dbNames = append(dbNames, name)
			}
		}
	}

	return dbNames, nil
}

func PersistentDatabaseName(databaseName string) string {
	return fmt.Sprintf(".%s_persistent", databaseName)
}
