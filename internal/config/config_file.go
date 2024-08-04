package config

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/go-errors/errors"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/viper"
)

const ConfigFileName = ".limerc"
const ConfigFileExt = "yaml"

func ReadConfigFile() {
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileExt)

	viper.AddConfigPath("/etc/limedb")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	homeDir := util.PanicIfErr(os.UserHomeDir())

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(errors.Errorf("Fatal error config file: %w", err))
		}

		slog.Warn("Could not find config file, creating new one", "Log code", "d44bb577")

		viper.SetDefault("limedb_home", filepath.Join(homeDir, ".limedb"))
		viper.SetDefault("log_mode", "stdout")
		viper.SetDefault("log_level", "info")
		viper.SetDefault("soft_deletion", true)
		viper.SetDefault("default_db", "")
		viper.SetDefault("remove_orphan_backups", true)

		if err := viper.WriteConfigAs(filepath.Join(homeDir, ConfigFileName+"."+ConfigFileExt)); err != nil {
			panic(errors.Errorf("Fatal error writing to config file: %w", err))
		}

	} else {
		slog.Info("Successfully read config file", "Log code", "b040b5d9")
	}
}
