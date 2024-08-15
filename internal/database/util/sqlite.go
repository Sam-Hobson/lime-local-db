package util

import (
	"github.com/huandu/go-sqlbuilder"
	"github.com/sam-hobson/internal/types"
	"github.com/sam-hobson/internal/util"
)

func CreateTableSql(tableName string, columns []*types.Column) (string, []interface{}) {
	util.Log("19529bb3").Info("Creating sqlite to create table.", "Table name", tableName, "Columns", columns)

	ctb := sqlbuilder.NewCreateTableBuilder()
	ctb.CreateTable(tableName).IfNotExists()

	for _, col := range columns {
		opts := make([]string, 10)
		opts = append(opts, col.Name)
		opts = append(opts, col.DataType.String())

		if col.DefaultVal != "" {
			opts = append(opts, "DEFAULT "+col.DefaultVal)
		}
		if col.NotNull {
			opts = append(opts, "NOT NULL")
		}
		if col.PrimaryKey {
			opts = append(opts, "PRIMARY KEY")
		}
		if col.AutoIncrememnt {
			opts = append(opts, "AUTOINCREMENT")
		}

		// if col.ForeignKey {
		// 	opts = append(opts, "FOREIGN KEY")
		// }

		ctb.Define(opts...)
	}

	return ctb.Build()
}

func InsertIntoTableSql(tableName string, entries map[string]string) (string, []interface{}) {
	util.Log("9391c009").Info("Creating sqlite to insert into table.", "Table name", tableName, "Entries", entries)

	ib := sqlbuilder.NewInsertBuilder()
	ib.InsertInto(tableName)

	keys := make([]string, len(entries))
	values := make([]interface{}, len(entries))
	var i = 0

	for key, value := range entries {
		keys[i] = key
		values[i] = value
		i++
	}

	ib.Cols(keys...)
	ib.Values(values...)

	return ib.Build()
}

func EntriesInTableWhereSql(tableName string, columns []string, args *sqlbuilder.Args, conditions ...string) (string, []interface{}) {
	util.Log("be384016").Info("Creating sqlite to check if entry exists in table.", "Table name", tableName, "Conditions", conditions)

	where := sqlbuilder.NewWhereClause()
	for _, cond := range conditions {
		where.AddWhereExpr(args, cond)
	}

	sb := sqlbuilder.NewSelectBuilder().Select(columns...).From(tableName).AddWhereClause(where)
	return sb.Build()
}
