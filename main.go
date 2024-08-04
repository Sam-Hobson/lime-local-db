package main

import (
	"log/slog"
	"os"

	"github.com/huandu/go-sqlbuilder"
	"github.com/sam-hobson/cmd"
	"github.com/sam-hobson/internal/config"
)

var (
	version = "beta"
	commit  = ""
)

var preConfigHandler = &slog.HandlerOptions{
	Level: slog.LevelError,
}

func main() {
	sqlbuilder.DefaultFlavor = sqlbuilder.SQLite

	// The default logger is set up prior to reading config, as reading config
	// contains info on how to set up logger further.
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, preConfigHandler)))

	config.ReadConfigFile()

	if err := cmd.NewCommand(version, commit).Execute(); err != nil {
		os.Exit(1)
	}
}
