package main

import (
	"os"

	"github.com/huandu/go-sqlbuilder"
	"github.com/sam-hobson/cmd"
	"github.com/sam-hobson/internal/config"
)

var (
	version = "beta"
	commit  = ""
)

func main() {
	sqlbuilder.DefaultFlavor = sqlbuilder.SQLite

	config.ReadConfigFile()

	if err := cmd.NewCommand(version, commit).Execute(); err != nil {
		os.Exit(1)
	}
}
