package removewhere

import (
	"log/slog"
	"strings"

	"github.com/go-errors/errors"
	"github.com/huandu/go-sqlbuilder"
	"github.com/sam-hobson/internal/database"
	"github.com/spf13/cobra"
)

// TODO: Improve command documentation
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm-entries-where ColumnName{Operation Value} ColumnName{Operation Value}...",
		Short:   "Remove entries conditionally from a database",
		Example: "limedb --db pets rm-entries-where name{like:%dog}",
		Args:    cobra.MinimumNArgs(1),

		RunE: run,
	}

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	where := sqlbuilder.NewWhereClause()

	for _, arg := range args {
		if err := parseColOperation(where, arg); err != nil {
			return err
		}
	}

	if rowsEffected, err := database.RemoveEntries(where); err != nil {
		return err
	} else {
		cmd.Printf("%d rows affected\n", rowsEffected)
		return nil
	}
}

func parseColOperation(where *sqlbuilder.WhereClause, arg string) error {
	colNameEnd := strings.Index(arg, ":")

	if (colNameEnd == -1) || (colNameEnd == 0) {
		slog.Error("Could not parse operation value.", "log_code", "5d769bb", "arg", arg)
		return errors.Errorf("Invalid rm-entries-where operation in %s", arg)
	}

	colName := arg[:colNameEnd]

	startCond := strings.Index(arg, "{")
	var op string

	if startCond == -1 {
		op = arg[colNameEnd+1:]
	} else {
		op = arg[colNameEnd+1 : startCond]
	}

	op = strings.ToUpper(strings.TrimSpace(op))
	sqlCond := sqlbuilder.NewCond()

	if op == "NULL" {
		where.AddWhereExpr(sqlCond.Args, sqlCond.IsNull(colName))
		return nil
	}
	if op == "NOTNULL" {
		where.AddWhereExpr(sqlCond.Args, sqlCond.IsNotNull(colName))
		return nil
	}

	endCond := strings.Index(arg, "}")

	if (startCond == -1) || (endCond == -1) {
		slog.Error("Could not parse operation value.", "log_code", "5f46ddba", "arg", arg)
		return errors.Errorf("Invalid rm-entries-where operation in %s", arg)
	}
	if (startCond == 0) || (endCond != len(arg)-1) {
		slog.Error("Could not parse operation value.", "log_code", "75378d9d", "arg", arg)
		return errors.Errorf("Invalid rm-entries-where operation in %s", arg)
	}

	opArg := arg[startCond+1 : endCond]

	switch op {
	case "=", "==":
		where.AddWhereExpr(sqlCond.Args, sqlCond.Equal(colName, opArg))
	case "!=", "<>":
		where.AddWhereExpr(sqlCond.Args, sqlCond.NotEqual(colName, opArg))
	case ">":
		where.AddWhereExpr(sqlCond.Args, sqlCond.LessThan(colName, opArg))
	case ">=":
		where.AddWhereExpr(sqlCond.Args, sqlCond.LessEqualThan(colName, opArg))
	case "<":
		where.AddWhereExpr(sqlCond.Args, sqlCond.GreaterThan(colName, opArg))
	case "<=":
		where.AddWhereExpr(sqlCond.Args, sqlCond.GreaterEqualThan(colName, opArg))
	case "LIKE":
		where.AddWhereExpr(sqlCond.Args, sqlCond.Like(colName, opArg))
	case "NOTLIKE":
		where.AddWhereExpr(sqlCond.Args, sqlCond.NotLike(colName, opArg))
	case "BETWEEN":
		if l, r, found := strings.Cut(opArg, ":"); !found {
			slog.Error("Could not parse operation.", "log_code", "040c7040", "arg", arg, "op", op)
			return errors.Errorf("Invalid rm-entries-where operation in %s", arg)
		} else {
			where.AddWhereExpr(sqlCond.Args, sqlCond.Between(colName, l, r))
		}
	case "NOTBETWEEN":
		if l, r, found := strings.Cut(opArg, ":"); !found {
			slog.Error("Could not parse operation.", "log_code", "98006b55", "arg", arg, "op", op)
			return errors.Errorf("Invalid rm-entries-where operation in %s", arg)
		} else {
			where.AddWhereExpr(sqlCond.Args, sqlCond.NotBetween(colName, l, r))
		}
	default:
		slog.Error("Could not parse operation.", "log_code", "75d191b9", "arg", arg, "op", op)
		return errors.Errorf("Invalid rm-entries-where operation in %s", arg)
	}

    return nil
}

// func parseColOperation(where *sqlbuilder.WhereClause, arg string) error {
// 	startCond := strings.Index(arg, "{")
// 	endCond := strings.Index(arg, "}")
//
// 	if (startCond == -1) || (endCond == -1) {
// 		slog.Error("Could not parse operation value.", "log_code", "04bd356", "arg", arg)
// 		return errors.Errorf("Invalid rm-entries-where operation in %s", arg)
// 	}
// 	if (startCond == 0) || (endCond != len(arg)-1) {
// 		slog.Error("Could not parse operation value.", "log_code", "d870b922", "arg", arg)
// 		return errors.Errorf("Invalid rm-entries-where operation in %s", arg)
// 	}
//
// 	colName := arg[:startCond]
// 	cond := arg[startCond+1 : endCond]
//
// 	sqlCond := sqlbuilder.NewCond()
//
// 	switch strings.ToUpper(strings.TrimSpace(cond)) {
// 	case "NULL":
// 		where.AddWhereExpr(sqlCond.Args, sqlCond.IsNull(colName))
// 		return nil
// 	case "NOTNULL":
// 		where.AddWhereExpr(sqlCond.Args, sqlCond.IsNotNull(colName))
// 		return nil
// 	}
//
// 	endOp := strings.Index(cond, ":")
//
// 	if (endOp == -1) || (endOp == 0) || (endOp == len(cond)-1) {
// 		slog.Error("Could not parse operation value.", "log_code", "17a7514a", "arg", arg)
// 		return errors.Errorf("Invalid rm-entries-where operation in %s", arg)
// 	}
//
// 	op := strings.TrimSpace(cond[:endOp])
// 	opArg := strings.TrimSpace(cond[endOp+1:])
//
// 	switch strings.ToUpper(op) {
// 	case "=", "==":
// 		where.AddWhereExpr(sqlCond.Args, sqlCond.Equal(colName, opArg))
// 	case "!=", "<>":
// 		where.AddWhereExpr(sqlCond.Args, sqlCond.NotEqual(colName, opArg))
// 	case ">":
// 		where.AddWhereExpr(sqlCond.Args, sqlCond.LessThan(colName, opArg))
// 	case ">=":
// 		where.AddWhereExpr(sqlCond.Args, sqlCond.LessEqualThan(colName, opArg))
// 	case "<":
// 		where.AddWhereExpr(sqlCond.Args, sqlCond.GreaterThan(colName, opArg))
// 	case "<=":
// 		where.AddWhereExpr(sqlCond.Args, sqlCond.GreaterEqualThan(colName, opArg))
// 	case "LIKE":
// 		where.AddWhereExpr(sqlCond.Args, sqlCond.Like(colName, opArg))
// 	case "NOTLIKE":
// 		where.AddWhereExpr(sqlCond.Args, sqlCond.NotLike(colName, opArg))
// 	case "BETWEEN":
// 		if l, r, found := strings.Cut(opArg, ":"); !found {
// 			slog.Error("Could not parse operation.", "log_code", "040c7040", "arg", arg, "op", op)
// 			return errors.Errorf("Invalid rm-entries-where operation in %s", arg)
// 		} else {
// 			where.AddWhereExpr(sqlCond.Args, sqlCond.Between(colName, l, r))
// 		}
// 	case "NOTBETWEEN":
// 		if l, r, found := strings.Cut(opArg, ":"); !found {
// 			slog.Error("Could not parse operation.", "log_code", "98006b55", "arg", arg, "op", op)
// 			return errors.Errorf("Invalid rm-entries-where operation in %s", arg)
// 		} else {
// 			where.AddWhereExpr(sqlCond.Args, sqlCond.NotBetween(colName, l, r))
// 		}
// 	default:
// 		slog.Error("Could not parse operation.", "log_code", "75d191b9", "arg", arg, "op", op)
// 		return errors.Errorf("Invalid rm-entries-where operation in %s", arg)
// 	}
//
// 	return nil
// }
