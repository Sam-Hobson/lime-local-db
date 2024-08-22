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

func RowsIntoSlice[T colType](rows *sql.Rows) []T {
	res := make([]T, 0)
	for rows.Next() {
		t := new(T)
		rows.Scan(t)
		res = append(res, *t)
	}
	return res
}

func RowsIntoSlice2[T colType, U colType](rows *sql.Rows) ([]T, []U) {
	res1, res2 := make([]T, 0), make([]U, 0)
	for rows.Next() {
		t, u := new(T), new(U)
		rows.Scan(t, u)
		res1, res2 = append(res1, *t), append(res2, *u)
	}
	return res1, res2
}

func RowsIntoSlice3[T colType, U colType, V colType](rows *sql.Rows) ([]T, []U, []V) {
	res1, res2, res3 := make([]T, 0), make([]U, 0), make([]V, 0)
	for rows.Next() {
		t, u, v := new(T), new(U), new(V)
		rows.Scan(t, u, v)
		res1, res2, res3 = append(res1, *t), append(res2, *u), append(res3, *v)
	}
	return res1, res2, res3
}

func RowsIntoSlice4[T colType, U colType, V colType, X colType](rows *sql.Rows) ([]T, []U, []V, []X) {
	res1, res2, res3, res4 := make([]T, 0), make([]U, 0), make([]V, 0), make([]X, 0)
	for rows.Next() {
		t, u, v, x := new(T), new(U), new(V), new(X)
		rows.Scan(t, u, v)
		res1, res2, res3, res4 = append(res1, *t), append(res2, *u), append(res3, *v), append(res4, *x)
	}
	return res1, res2, res3, res4
}
