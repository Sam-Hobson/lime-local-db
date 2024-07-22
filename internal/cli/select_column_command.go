package cli

import (
	"github.com/sam-hobson/internal/operations"
	"github.com/spf13/cobra"
)

const selColumnCommandUsage = `with-col ColName`
const selColumnShortCommandUsage = `col ColName`
const selColumnCommandExample = `limedb with-col colName`
const selColumnShortCommandExample = `limedb col colName`
const selColumnCommandShort = `Select a column to use for operations.`

var SelColumnCommand = &cobra.Command{
	Use:     selColumnCommandUsage,
	Example: selColumnCommandExample,
	Short:   selDbCommandShort,
	Args:    cobra.ExactArgs(1),

	RunE: processColumnDbCmd,
}

var SelColumnShortCommand = &cobra.Command{
	Use:     selDbShortCommandUsage,
	Example: selDbShortCommandExample,
	Short:   selDbCommandShort,
	Args:    cobra.ExactArgs(1),

	RunE: processColumnDbCmd,
}

func processColumnDbCmd(_ *cobra.Command, args []string) error {
	operations.SelectCol(args[0])
	return nil
}
