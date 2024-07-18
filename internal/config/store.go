package config

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	LimeDir = ".limedb"
)

func createRootDir() error {
	rootPath := filepath.Join(home, LimeDir)
	err := os.MkdirAll(rootPath, os.ModePerm)
	return err
}

func GetRelDir(relPath string) (fs.FS, error) {
	path := filepath.Join(home, LimeDir, relPath)
	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		slog.Error("Could not create directory.", "Hash", "791e81d6", "Directory", path)
		return nil, err
	}

	slog.Info("Created/retrieved dir.", "Hash", "17a6838b", "path", path)
	return os.DirFS(path), nil
}
