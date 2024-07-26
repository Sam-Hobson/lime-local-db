package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/go-errors/errors"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/viper"
)

const ConfigFileName = "limerc"
const ConfigFileExt = "yaml"

func ReadConfigFile() {
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileExt)

	viper.AddConfigPath("/etc/limedb")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

    homeDir := util.PanicIfErr(os.UserHomeDir())

    viper.SetDefault("limedbHome", filepath.Join(homeDir, ".limedb"))
    viper.SetDefault("softDeletion", true)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(errors.Errorf("Fatal error config file: %w", err))
		}

        slog.Warn("Could not find config file, creating new one", "log_code", "d44bb577")

		err = viper.WriteConfigAs(fmt.Sprintf("%s/%s.%s", homeDir, ConfigFileName, ConfigFileExt))

		if err != nil {
			panic(errors.Errorf("Fatal error creating config file: %w", err))
		}
	} else {
        slog.Info("Successfully read config file", "log_code", "b040b5d9")
    }

}
