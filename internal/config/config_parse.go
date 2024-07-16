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
	LimeDir    = ".limedb"
)

func ParseConfig(homeDir string) (*cp.ConfigParser, error) {
	configPath := filepath.Join(homeDir, ConfigName)

	parser, err := cp.NewConfigParserFromFile(configPath)

	if err != nil {
		slog.Error("Could not open config file.\n", "path", configPath)
		return nil, err
	}

	slog.Info("Successfully parsed config file.")

	return parser, nil
}

func CreateDefaultConfig(homeDir string, rootDir string) error {
	configPath := filepath.Join(homeDir, ConfigName)

	template := getDefaultTemplate(rootDir)

	slog.Info("Creating default config.",
		"path", configPath,
		"homeDir", homeDir,
		"rootDir", rootDir)

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

	slog.Info("Successfully created default config.\n",
		"path", configPath,
		"homeDir", homeDir,
		"rootDir", rootDir)

	return nil
}

const defaultTemplate = `[DEBUG]
log_location = {LOG_LOCATION}

[STORE]
root_dir = {ROOT_DIR}
`

func getDefaultTemplate(rootDir string) string {
	var template = defaultTemplate

	defaultMap := map[string]string{
		"{LOG_LOCATION}": "",
		"{ROOT_DIR}":     rootDir,
	}

	// Yes, this is horribly inefficient but it will do for now
	for k, v := range defaultMap {
		template = strings.ReplaceAll(template, k, v)
	}

	return template
}
