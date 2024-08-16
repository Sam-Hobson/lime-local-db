package addentry

import (
	"fmt"
	"strings"

	"github.com/go-errors/errors"
	"github.com/sam-hobson/internal/database"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
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

	util.Log("7a8f5e35").Info("Parsed add-entry arguments.", "Args", entryValues)

	databaseName := state.ApplicationState().GetSelectedDb()
	if databaseName == "" {
		util.Log("0562d5bb").Error("Cannot add entry if not database is specified.")
		return errors.Errorf("Cannot add entry if not database is specified")
	}

	if err := database.AddEntry(databaseName, entryValues); err != nil {
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Successfully added entry")
	return nil
}

func parseColValue(value string) (string, string, error) {
	startVal := strings.Index(value, "{")
	endVal := strings.Index(value, "}")

	if (startVal == -1) || (endVal == -1) {
		util.Log("1361631e").Error("Could not parse entry value.", "Value", value)
		return "", "", errors.Errorf("Invalid add-entry value in %s, value not present", value)
	}
	if (startVal == 0) || (endVal != len(value)-1) {
		util.Log("f6fa2b05").Error("Could not parse entry value.", "Value", value)
		return "", "", errors.Errorf("Invalid add-entry value in %s, malformed input", value)
	}

	return value[:startVal], value[startVal+1 : endVal], nil
}
