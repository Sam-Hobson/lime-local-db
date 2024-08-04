package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/sam-hobson/internal/util"
	"github.com/spf13/viper"
)

func GetConfigLogWriter() *os.File {
	var writer *os.File

	switch strings.ToUpper(viper.GetString("log_mode")) {
	case "STDOUT":
		writer = os.Stdout
	case "STDERR":
		writer = os.Stderr
	case "FILE":
		relFs := util.NewRelativeFsManager(viper.GetString("limedb_home"))
		if w, err := relFs.OpenFile("", "limedb.log"); err != nil {
			fatal("Unable open/create log file at %s", filepath.Join(viper.GetString("limedb_home"), "limedb.log"))
		} else {
			writer = w
		}
	default:
		fatal("Invalid log_mode in config file. Not one of \"stdout\", \"stderr\", \"file\".")
	}

	return writer
}

func GetConfigLogLevel() slog.Level {
	switch strings.ToUpper(viper.GetString("log_level")) {
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "DEBUG":
		return slog.LevelDebug
	case "ERROR":
		return slog.LevelError
	default:
		fatal("Invalid log_level in config file. Not one of \"info\", \"warn\", \"debug\", \"error\".")
	}

	return slog.LevelInfo
}

func fatal(s string, args ...any) {
	fmt.Printf(s+"\n", args...)
	os.Exit(1)
}
