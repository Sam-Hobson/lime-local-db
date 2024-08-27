package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/sam-hobson/cmd"
	"github.com/sam-hobson/internal/config"
	"github.com/sam-hobson/internal/util"
)

var (
	version = "beta"
	commit  = ""
)

var preConfigHandler = &slog.HandlerOptions{
	Level: slog.LevelError,
}

func main() {
	util.SetSessionId(time.Now().UnixMicro())
	sqlbuilder.DefaultFlavor = sqlbuilder.SQLite

	// The default logger is set up prior to reading config, as reading config
	// contains info on how to set up logger further.
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, preConfigHandler)))

	config.ReadConfigFile()

	if err := cmd.NewCommand(version, commit).Execute(); err != nil {
		panic(err)
	}
}
