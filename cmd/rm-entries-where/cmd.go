package removewhere

import (
	"strings"

	"github.com/go-errors/errors"
	"github.com/huandu/go-sqlbuilder"
	"github.com/sam-hobson/internal/database"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
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

	databaseName := state.ApplicationState().GetSelectedDb()
	if databaseName == "" {
		util.Log("53148429").Error("Cannot remove entries if not database is specified.")
		return errors.Errorf("Cannot remove entries if not database is specified")
	}

	if rowsEffected, err := database.RemoveEntries(databaseName, where); err != nil {
		return err
	} else {
		cmd.Printf("%d rows affected\n", rowsEffected)
		return nil
	}
}

func parseColOperation(where *sqlbuilder.WhereClause, arg string) error {
	colNameEnd := strings.Index(arg, ":")

	if (colNameEnd == -1) || (colNameEnd == 0) {
		util.Log("5d769bb").Error("Could not parse operation value.", "Arg", arg)
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
		util.Log("5f46ddba").Error("Could not parse operation value.", "Arg", arg)
		return errors.Errorf("Invalid rm-entries-where operation in %s", arg)
	}
	if (startCond == 0) || (endCond != len(arg)-1) {
		util.Log("75378d9d").Error("Could not parse operation value.", "Arg", arg)
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
			util.Log("040c7040").Error("Could not parse operation.", "Arg", arg, "Op", op)
			return errors.Errorf("Invalid rm-entries-where operation in %s", arg)
		} else {
			where.AddWhereExpr(sqlCond.Args, sqlCond.Between(colName, l, r))
		}
	case "NOTBETWEEN":
		if l, r, found := strings.Cut(opArg, ":"); !found {
			util.Log("98006b55").Error("Could not parse operation.", "Arg", arg, "Op", op)
			return errors.Errorf("Invalid rm-entries-where operation in %s", arg)
		} else {
			where.AddWhereExpr(sqlCond.Args, sqlCond.NotBetween(colName, l, r))
		}
	default:
		util.Log("75d191b9").Error("Could not parse operation.", "Arg", arg, "Op", op)
		return errors.Errorf("Invalid rm-entries-where operation in %s", arg)
	}

	return nil
}
