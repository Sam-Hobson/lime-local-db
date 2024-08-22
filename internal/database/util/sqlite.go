package util

import (
	"database/sql"

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

type colType interface {
	int | float32 | float64 | string
}

func RowsIntoSlice[C colType](rows *sql.Rows) []C {
	res := make([]C, 0)
	for rows.Next() {
		c := new(C)
		rows.Scan(c)
		res = append(res, *c)
	}
	return res
}

func RowsIntoSlice2[C1 colType, C2 colType](rows *sql.Rows) ([]C1, []C2) {
	res1, res2 := make([]C1, 0), make([]C2, 0)
	for rows.Next() {
		c1, c2 := new(C1), new(C2)
		rows.Scan(c1, c2)
		res1, res2 = append(res1, *c1), append(res2, *c2)
	}
	return res1, res2
}

func RowsIntoSlice3[C1 colType, C2 colType, C3 colType](rows *sql.Rows) ([]C1, []C2, []C3) {
	res1, res2, res3 := make([]C1, 0), make([]C2, 0), make([]C3, 0)
	for rows.Next() {
		c1, c2, c3 := new(C1), new(C2), new(C3)
		rows.Scan(c1, c2, c3)
		res1, res2, res3 = append(res1, *c1), append(res2, *c2), append(res3, *c3)
	}
	return res1, res2, res3
}

func RowsIntoSlice4[C1 colType, C2 colType, C3 colType, C4 colType](rows *sql.Rows) ([]C1, []C2, []C3, []C4) {
	res1, res2, res3, res4 := make([]C1, 0), make([]C2, 0), make([]C3, 0), make([]C4, 0)
	for rows.Next() {
		c1, c2, c3, c4 := new(C1), new(C2), new(C3), new(C4)
		rows.Scan(c1, c2, c3)
		res1, res2, res3, res4 = append(res1, *c1), append(res2, *c2), append(res3, *c3), append(res4, *c4)
	}
	return res1, res2, res3, res4
}
