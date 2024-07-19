package cli

import (
	"log/slog"
	"strings"

	"github.com/go-errors/errors"
	op "github.com/sam-hobson/internal/operations"
	"github.com/spf13/cobra"
)

var NewDbCommand = &cobra.Command{
	Use:     "new-db [Database name] \"[[Database columns]]\"",
	Example: "limedb new-db petdb \"[name gender breed]\"",
	Short:   "Create a new database",
    Args: cobra.ExactArgs(2),

	DisableFlagsInUseLine: true, // TODO: Remove this when flags are added

	RunE: processNewDbCmd,
}

func processNewDbCmd(cmd *cobra.Command, args []string) error {
	// TODO: Add support for various flags
	anyFlagSet := cmd.Flags().NFlag() != 0

	if anyFlagSet {
		slog.Error("new-db does not take any flags.", "log_code", "998073c6")
		return errors.Errorf("new-db does not take any flags.")
	}

	name := args[0]
	attrs := strings.Fields(args[1])

	err := op.NewDb(name, attrs)

	if err != nil {
		slog.Error("setup command failed failed.", "log_code", "0f36afa0")
		return err
	}

	return nil
}
