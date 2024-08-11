package db

import (
	"github.com/huandu/go-sqlbuilder"
	"github.com/sam-hobson/internal/database/masterdatabase"
	"github.com/sam-hobson/internal/state"
	"github.com/spf13/cobra"
)

func lsDbCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "ls",
		Short:     "List databases",
		Example:   "limedb db ls",
		Args:      cobra.NoArgs,
		ValidArgs: dbNames(),

		RunE: runLsDbCommand,
	}

	cmd.Flags().Bool("all", false, "Show all details about databases")

	return cmd
}

func runLsDbCommand(cmd *cobra.Command, _ []string) error {
    databaseName := state.ApplicationState().GetSelectedDb()

	cond := sqlbuilder.NewCond()
	where := sqlbuilder.NewWhereClause()

    if databaseName != "" {
        where.AddWhereExpr(cond.Args, cond.Equal("name", databaseName))
    }
    where.AddWhereExpr(cond.Args, cond.IsNull("softdeleted"))
    where.AddWhereExpr(cond.Args, cond.IsNull("harddeleted"))

    res, err := masterdatabase.QueryTables(where, "name")
    if err != nil {
        return err
    }
    defer res.Close()

    for res.Next() {
        var name string
        res.Scan(&name)
        cmd.Println(name)
    }

    return nil
}
