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
		opts := make([]string, 0, 10)
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

		ctb.Define(opts...)
	}
	for _, col := range columns {
		if col.ForeignKey == nil {
			continue
		}

		opts := make([]string, 0, 10)
		opts = append(opts, "FOREIGN KEY")
		opts = append(opts, "("+col.Name+")")
		opts = append(opts, "REFERENCES")
		opts = append(opts, col.ForeignKey.Table)
		opts = append(opts, "("+col.ForeignKey.Col+")")

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

func RowsIntoSlice[C colType](rows *sql.Rows) ([]C, error) {
	res := make([]C, 0)
	for rows.Next() {
		c := new(C)
		if err := rows.Scan(c); err != nil {
			util.Log("835720e2").Warn("Error while reading rows into slice.")
			return nil, err
		}
		res = append(res, *c)
	}
	return res, nil
}

func RowsIntoSlice2[C1 colType, C2 colType](rows *sql.Rows) ([]C1, []C2, error) {
	res1, res2 := make([]C1, 0), make([]C2, 0)
	for rows.Next() {
		c1, c2 := new(C1), new(C2)
		if err := rows.Scan(c1, c2); err != nil {
			util.Log("13903e2e").Warn("Error while reading rows into slice.")
			return nil, nil, err
		}
		res1, res2 = append(res1, *c1), append(res2, *c2)
	}
	return res1, res2, nil
}

func RowsIntoSlice3[C1 colType, C2 colType, C3 colType](rows *sql.Rows) ([]C1, []C2, []C3, error) {
	res1, res2, res3 := make([]C1, 0), make([]C2, 0), make([]C3, 0)
	for rows.Next() {
		c1, c2, c3 := new(C1), new(C2), new(C3)
		if err := rows.Scan(c1, c2, c3); err != nil {
			util.Log("73743de3").Warn("Error while reading rows into slice.")
			return nil, nil, nil, err
		}
		res1, res2, res3 = append(res1, *c1), append(res2, *c2), append(res3, *c3)
	}
	return res1, res2, res3, nil
}

func RowsIntoSlice4[C1 colType, C2 colType, C3 colType, C4 colType](rows *sql.Rows) ([]C1, []C2, []C3, []C4, error) {
	res1, res2, res3, res4 := make([]C1, 0), make([]C2, 0), make([]C3, 0), make([]C4, 0)
	for rows.Next() {
		c1, c2, c3, c4 := new(C1), new(C2), new(C3), new(C4)
		if err := rows.Scan(c1, c2, c3, c4); err != nil {
			util.Log("73743de3").Warn("Error while reading rows into slice.")
			return nil, nil, nil, nil, err
		}
		res1, res2, res3, res4 = append(res1, *c1), append(res2, *c2), append(res3, *c3), append(res4, *c4)
	}
	return res1, res2, res3, res4, nil
}
