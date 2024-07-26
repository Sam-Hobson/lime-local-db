package cmd

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/go-errors/errors"
	newdb "github.com/sam-hobson/cmd/new-db"
	rmdb "github.com/sam-hobson/cmd/rm-db"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCommand(version, commit string) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "limedb",
		Short:             "Application for interacting with a simple database",
		Long:              "TODO: This",
		Version:           buildVersion(version, commit),
		PersistentPreRunE: preRun,
	}

	cmd.PersistentFlags().StringSlice("with-config", nil, "Override a configuration option during the execution of this command.")

	cmd.AddCommand(newdb.NewCommand(), rmdb.NewCommand())

	return cmd
}

func preRun(cmd *cobra.Command, args []string) error {
	configChanges := util.PanicIfErr(cmd.Flags().GetStringSlice("with-config"))

	if configChanges != nil {
		slog.Info("--with-config provided.", "log_code", "4d720ec4", "Config-changes", fmt.Sprintf("%v", configChanges))

		for _, change := range configChanges {
			key, value, found := strings.Cut(change, ":")

			if !found {
				slog.Error("--with-config flag has malformed input.", "log_code", "d05c9983", "Changes", change)
				return errors.Errorf("--with-config flag has malformed input.")
			}

			viper.Set(key, value)
		}
	}

	return nil
}

func buildVersion(version, commit string) string {
	return fmt.Sprintf("%s(%s)", version, commit)
}
