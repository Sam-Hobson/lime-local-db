package cli

import (
	"github.com/sam-hobson/internal/operations"
	"github.com/spf13/cobra"
)

const selDbCommandUsage = `with-db DatabaseName`
const selDbShortCommandUsage = `db DatabaseName`
const selDbCommandExample = `limedb with-db petdb`
const selDbShortCommandExample = `limedb db petdb`
const selDbCommandShort = `Select a database to use for operations.`

var SelDbCommand = &cobra.Command{
	Use:     selDbCommandUsage,
	Example: selDbCommandExample,
	Short:   selDbCommandShort,
	Args:    cobra.ExactArgs(1),

	RunE: processSelDbCmd,
}

var SelDbShortCommand = &cobra.Command{
	Use:     selDbShortCommandUsage,
	Example: selDbShortCommandExample,
	Short:   selDbCommandShort,
	Args:    cobra.ExactArgs(1),

	RunE: processSelDbCmd,
}

func processSelDbCmd(_ *cobra.Command, args []string) error {
	operations.SelectDb(args[0])
	return nil
}

func init() {
	SelDbCommand.AddCommand(SelColumnCommand)
}
