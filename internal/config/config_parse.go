package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	cp "github.com/bigkevmcd/go-configparser"
	"github.com/go-errors/errors"
)

const (
	ConfigName = ".limerc"
	LimeDir    = ".limedb"
)

var home = os.Getenv("HOME")

func ParseConfig() (*cp.ConfigParser, error) {
	if home == "" {
		slog.Error("Cannot create config file as $HOME is undefined")
		return nil, errors.New("$HOME is undefined")
	}

	configPath := filepath.Join(home, ConfigName)

	// If the config file hasn't yet been created...
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		err := createDefaultConfigFile(configPath)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	parser, err := cp.NewConfigParserFromFile(configPath)

	if err != nil {
		slog.Error("Could not open config file.\n", "path", configPath)
		return nil, err
	}

	slog.Info("Successfully parsed config file.")

	return parser, nil
}

func createDefaultConfigFile(configPath string) error {
	var template = defaultTemplate

	// Yes, this is horribly inefficient but it will do for now
	for k, v := range defaultMapping {
		template = strings.ReplaceAll(template, k, v)
	}

	slog.Info("Creating default config.", "path", configPath)

	file, err := os.OpenFile(configPath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	defer file.Close()

	if err != nil {
		slog.Error("Cannot create config file.\n", "path", configPath)
		return err
	}

	_, err = file.WriteString(template)

	if err != nil {
		slog.Error("Cannot create config file.\n", "path", configPath)
		return err
	}

	slog.Info("Successfully created default config.\n", "path", configPath)
	return nil
}

const defaultTemplate = `[DEBUG]
log_location = {LOG_LOCATION}

[STORE]
store_location = {STORE_LOCATION}
`

var defaultMapping = map[string]string{
	"{LOG_LOCATION}":   "",
	"{STORE_LOCATION}": filepath.Join(home, LimeDir, "store"),
}
