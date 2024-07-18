package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	cp "github.com/bigkevmcd/go-configparser"
)

const (
	ConfigName = ".limerc"
)

const (
	StoreSection = "STORE"
	DebugSection = "DEBUG"
)

var home = os.Getenv("HOME")

var config *cp.ConfigParser
var configParsed bool

func GetConfig() *cp.ConfigParser {
	if !configParsed {
		return nil
	}

	return config
}

func ConfigExists() bool {
	configPath := filepath.Join(home, ConfigName)
	_, err := os.Stat(configPath)
	return err == nil
}

func ParseConfig() error {
	configPath := filepath.Join(home, ConfigName)

	parser, err := cp.NewConfigParserFromFile(configPath)

	if err != nil {
		slog.Error("Could not open config file.\n", "Hash", "10e671cd", "path", configPath)
		return err
	}

	slog.Info("Successfully parsed config file.", "Hash", "caf29dea")

	configParsed = true
	config = parser

	return nil
}

func CreateDefaultConfig() error {
	configPath := filepath.Join(home, ConfigName)

	template := getDefaultTemplate(home)

	slog.Info("Creating default config.",
        "Hash", "cb196a26",
		"path", configPath,
		"homeDir", home)

	file, err := os.OpenFile(configPath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	defer file.Close()

	if err != nil {
		slog.Error("Cannot create config file.\n", "Hash", "3f93600b", "path", configPath)
		return err
	}

	_, err = file.WriteString(template)

	if err != nil {
		slog.Error("Cannot create config file.\n", "Hash", "b790cf52", "path", configPath)
		return err
	}

	slog.Info("Successfully created default config.\n",
        "Hash", "8ea69957",
		"path", configPath,
		"homeDir", home)

	return nil
}

const defaultTemplate = `[DEBUG]
log_type = stdout
log_location = {LOG_LOCATION}

[STORE]
root_dir = {ROOT_DIR}
`

func getDefaultTemplate(homeDir string) string {
	var template = defaultTemplate

	defaultMap := map[string]string{
		"{LOG_LOCATION}": "",
		"{ROOT_DIR}":     filepath.Join(homeDir, LimeDir),
	}

	// Yes, this is horribly inefficient but it will do for now
	for k, v := range defaultMap {
		template = strings.ReplaceAll(template, k, v)
	}

	return template
}
