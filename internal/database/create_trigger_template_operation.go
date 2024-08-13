package database

import (
	"fmt"
	"slices"

	"github.com/go-errors/errors"
	dbutil "github.com/sam-hobson/internal/database/util"
	"github.com/sam-hobson/internal/util"
)

const createTriggerTemplate = `CREATE TRIGGER %s
   BEFORE INSERT ON %s
BEGIN
END;
`

func TriggerTemplate(selectedDatabase, triggerName string) (string, error) {
	util.Log("fd238bcc").Info("Beginning trigger template operation.", "Selected database", selectedDatabase, "Trigger name", triggerName)

	if exists, err := dbutil.SqliteDatabaseExists(selectedDatabase); err != nil {
		util.Log("8e328019").Error("Cannot check if database exists, cannot create trigger template.", "Selected database", selectedDatabase)
		return "", err
	} else if !exists {
		util.Log("5151b742").Error("Cannot create trigger template if database doesn't exist.", "Selected database", selectedDatabase)
		return "", errors.Errorf("Cannot create trigger template if database doesn't exist")
	}

	if triggers, err := dbutil.DefinedTriggers(selectedDatabase); err != nil {
		util.Log("460e1305").Error("Cannot get defined triggers on table.", "Selected database", selectedDatabase)
		return "", err
	} else if slices.IndexFunc(triggers, func(trigger *dbutil.Trigger) bool { return trigger.Name == triggerName }) != -1 {
        util.Log("467910b2").Error("Cannot create trigger template as trigger already exists.", "Selected database", selectedDatabase, "Trigger name", triggerName)
        return "", errors.Errorf("Cannot create trigger template as trigger already exists")
	}

	res := fmt.Sprintf(createTriggerTemplate, triggerName, selectedDatabase)
	util.Log("a3fd4c31").Info("Successfully completed trigger template operation.", "Template", res)

	return res, nil
}
