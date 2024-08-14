package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/go-errors/errors"
	addentry "github.com/sam-hobson/cmd/add-entry"
	"github.com/sam-hobson/cmd/backup"
	"github.com/sam-hobson/cmd/db"
	rmentriesall "github.com/sam-hobson/cmd/rm-entries-all"
	rmentrieswhere "github.com/sam-hobson/cmd/rm-entries-where"
	"github.com/sam-hobson/cmd/trigger"
	"github.com/sam-hobson/internal/config"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logWriter *os.File

func NewCommand(version, commit string) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "limedb",
		Short:             "Application for interacting with a simple database",
		Long:              "TODO: This",
		Version:           buildVersion(version, commit),
		PersistentPreRunE: preRun,
		PersistentPostRun: postRun,
	}

	cmd.PersistentFlags().StringSlice("with-config", nil, "Override a configuration option during the execution of this command.")
	cmd.PersistentFlags().StringP("db", "d", "", "Choose the database to perform operations on.")

	cmd.AddCommand(
		db.NewCommand(),
		addentry.NewCommand(),
		rmentriesall.NewCommand(),
		rmentrieswhere.NewCommand(),
		backup.NewCommand(),
        trigger.NewCommand(),
	)

    // Not sure why this is necessary, but output goes to stderr by default if this isn't here
    cmd.SetOut(os.Stdout)

	return cmd
}

func preRun(cmd *cobra.Command, args []string) error {
	configChanges := util.PanicIfErr(cmd.Flags().GetStringSlice("with-config"))

	// Set single run config changes
	if len(configChanges) != 0 {
		for _, change := range configChanges {
			key, value, found := strings.Cut(change, ":")

			if !found {
				return errors.Errorf("--with-config flag has malformed input.")
			}

			viper.Set(key, value)
		}
	}

	// Set the logger up based on config
	logWriter = config.GetConfigLogWriter()
	level := config.GetConfigLogLevel()
	handler := slog.NewJSONHandler(logWriter, &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(handler))

	util.Log("09b726e0").Info("New limedb session starting.")

	if len(configChanges) != 0 {
		util.Log("3943d45c").Info("Successfully applied --with-config configuration.", "Config changes", configChanges)
	}

	if db := viper.GetString("default_db"); db != "" {
		state.ApplicationState().SetSelectedDb(db)
	}
	if limedbHome := viper.GetString("limedb_home"); limedbHome != "" {
		state.ApplicationState().SetLimedbHome(limedbHome)
	}

	// Process flags
	if selectedDb := util.PanicIfErr(cmd.Flags().GetString("db")); selectedDb != "" {
		state.ApplicationState().SetSelectedDb(selectedDb)
	}

	return nil
}

func postRun(cmd *cobra.Command, args []string) {
	logWriter.Close()
}

func buildVersion(version, commit string) string {
	return fmt.Sprintf("%s(%s)", version, commit)
}
