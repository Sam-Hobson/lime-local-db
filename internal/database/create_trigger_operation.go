package database

import (
	"fmt"
	"slices"
	"time"

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

func CreateTrigger(databaseName string, triggerName string, ttype types.TriggerType, body string, comment string) error {
	util.Log("e805a489").Info("Beginning create trigger operation.",
		"Database name", databaseName,
		"Trigger name", triggerName,
		"Trigger type", ttype,
		"Body", body)

	triggerStr, err := TriggerTemplate(databaseName, triggerName, ttype.String(), body)
	if err != nil {
		return err
	}

	return CreateTriggerRaw(databaseName, triggerName, triggerStr, comment)
}

func CreateTriggerRaw(databaseName, triggerName, triggerStr, comment string) error {
	util.Log("79915c34").Info("Creating trigger from string operation.", "Database name", databaseName, "Trigger string", triggerStr)

	db, err := dbutil.OpenSqliteDatabaseIfExists(databaseName)
	if err != nil {
		return err
	}
	defer db.Close()

	if exists, err := dbutil.TriggerNameExists(db, triggerName); err != nil {
		return err
	} else if exists {
		util.Log("add8b5cf").Error("Cannot create trigger as a trigger by name already exists.", "Database name", databaseName, "Trigger name", triggerName)
		return errors.Errorf("Cannot create trigger as a trigger by provided name already exists")
	}

	insertStr, args := dbutil.InsertIntoTableSql("triggers", map[string]string{
		"name":         triggerName,
		"date_created": time.Now().Format(time.RFC3339),
		"comment":      comment,
	})
	if _, err := db.Exec(insertStr, args...); err != nil {
		util.Log("465dca8d").Error("Could not create entry in triggers table.", "Database name", databaseName, "Trigger name", triggerName)
		return errors.Errorf("Could not create entry in triggers table")
	}

	// TODO: use the rowid returned to create a foreign key in triggers referring to sqlite_master.
	if _, err := ExecRawSql(databaseName, triggerStr); err != nil {
		return err
	}

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
