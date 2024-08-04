package addentry

import (
	"log/slog"
	"strings"

	"github.com/go-errors/errors"
	"github.com/sam-hobson/internal/database"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add-entry ColumnName{Value} ColumnName{Value}...",
		Short:   "Add a new entry to a database.",
		Example: "limedb --db pets add-entry name{John} gender{M} age{5}",
		Args:    cobra.MinimumNArgs(1),

		RunE: run,
	}

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	entryValues := make(map[string]string)

	for _, value := range args {
		if key, value, err := parseColValue(value); err != nil {
			return err
		} else {
			entryValues[key] = value
		}
	}

	slog.Info("Parsed add-entry arguments.", "Log code", "7a8f5e35", "Args", entryValues)
    return database.AddEntry(entryValues)
}

func parseColValue(value string) (string, string, error) {
	startVal := strings.Index(value, "{")
	endVal := strings.Index(value, "}")

	if (startVal == -1) || (endVal == -1) {
		slog.Error("Could not parse entry value.", "Log code", "1361631e", "Value", value)
		return "", "", errors.Errorf("Invalid add-entry value in %s, value not present", value)
	}
	if (startVal == 0) || (endVal != len(value)-1) {
		slog.Error("Could not parse entry value.", "Log code", "f6fa2b05", "Value", value)
		return "", "", errors.Errorf("Invalid add-entry value in %s, malformed input", value)
	}

	return value[:startVal], value[startVal+1 : endVal], nil
}
