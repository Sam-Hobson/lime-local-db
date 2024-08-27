package database

import (
	"github.com/huandu/go-sqlbuilder"
	"github.com/sam-hobson/internal/types"
)

var backupColumns = []*types.Column{
	{
		Name:     "date",
		DataType: types.ColumnTextDataType,
		NotNull:  true,
	},
	{
		Name:       "backup_name",
		DataType:   types.ColumnTextDataType,
		NotNull:    true,
		PrimaryKey: true,
	},
	{
		Name:     "comment",
		DataType: types.ColumnTextDataType,
	},
}

type Backup struct {
	Name    string `db:"backup_name" fieldtag:"pk"`
	Date    string `db:"date"`
	Comment string `db:"comment"`
}

var BackupStruct = sqlbuilder.NewStruct(new(Backup))

var triggerColumns = []*types.Column{
	{
		Name:       "sqlite_master_rowid",
		DataType:   types.ColumnIntDataType,
		NotNull:    true,
		PrimaryKey: true,
		ForeignKey: &types.ForeignKey{
			Table: "sqlite_master",
			Col:   "rowid",
		},
	},
	{
		Name:     "date_created",
		DataType: types.ColumnTextDataType,
		NotNull:  true,
	},
	{
		Name:     "trigger_type",
		DataType: types.ColumnTextDataType,
	},
	{
		Name:     "comment",
		DataType: types.ColumnTextDataType,
	},
}

type Trigger struct {
	SqliteMasterRowid int    `db:"sqlite_master_rowid" fieldtag:"pk"`
	Date              string `db:"date_created"`
	TriggerType       string `db:"trigger_type"`
	Comment           string `db:"comment"`
}

var TriggerStruct = sqlbuilder.NewStruct(new(Trigger))
