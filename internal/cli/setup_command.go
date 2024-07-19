package cli

import (
	"log/slog"

	"github.com/go-errors/errors"
	op "github.com/sam-hobson/internal/operations"
	"github.com/spf13/cobra"
)

var SetupCommand = &cobra.Command{
	Use:   "setup",
	Short: "Create a default configuration file at ~/.limerc",
	Args:  cobra.ExactArgs(0),

	RunE: processSetupCmd,
}

func processSetupCmd(cmd *cobra.Command, _ []string) error {
	anyFlagSet := cmd.Flags().NFlag() != 0

	if anyFlagSet {
		slog.Error("setup does not take flags.", "log_code", "63203f93")
		return errors.Errorf("setup does not take flags")
	}

	err := op.Setup()

	if err != nil {
		slog.Error("setup command failed failed.", "log_code", "0f36afa0")
		return err
	}

	return nil
}
