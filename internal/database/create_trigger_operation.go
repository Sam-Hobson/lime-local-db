package database

import (
	"fmt"
	"slices"

	"github.com/go-errors/errors"
	dbutil "github.com/sam-hobson/internal/database/util"
	"github.com/sam-hobson/internal/types"
	"github.com/sam-hobson/internal/util"
)

const triggerTemplate = `CREATE TRIGGER %s
   %s ON %s
BEGIN
    %s
END;
`

func CreateTrigger(databaseName string, triggerName string, ttype types.TriggerType, body string) error {
	util.Log("e805a489").Info("Beginning create trigger operation.",
		"Database name", databaseName,
		"Trigger name", triggerName,
		"Trigger type", ttype,
		"Body", body)

	triggerStr, err := TriggerTemplate(databaseName, triggerName, ttype.String(), body)
	if err != nil {
		return err
	}

	return CreateTriggerRaw(databaseName, triggerName, triggerStr)
}

func CreateTriggerRaw(databaseName, triggerName, triggerStr string) error {
	util.Log("79915c34").Info("Beginning create trigger from string operation.", "Database name", databaseName, "Trigger string", triggerStr)

    if _, err := ExecRawSql(databaseName, triggerStr); err != nil {
        return err
    }

	util.Log("27afb04c").Info("Successfully created trigger from string.", "Database name", databaseName)
	return nil
}

func TriggerTemplate(selectedDatabase, triggerName, ttype, body string) (string, error) {
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

	res := fmt.Sprintf(triggerTemplate, triggerName, ttype, selectedDatabase, body)
	util.Log("a3fd4c31").Info("Successfully completed trigger template operation.", "Template", res)

	return res, nil
}
