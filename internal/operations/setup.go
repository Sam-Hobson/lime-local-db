package operations

import (
	"log/slog"

	"github.com/sam-hobson/internal/config"
)

func Setup() error {
	slog.Info("Beginning setup operation.", "log_code", "8b6822bb")

	if config.ConfigExists() {
		slog.Warn("Setup called when there already exists a config file. Nothing done.", "log_code", "d3a9e666")
		return nil
	}

	err := config.CreateDefaultConfig()
	if err != nil {
		return err
	}

	err = config.ParseConfig()
	return err
}
