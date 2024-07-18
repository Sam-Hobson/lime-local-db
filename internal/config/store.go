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

func CreateFile(relPath string, fileName string) error {
	err := CreateDir(FullPath(relPath))

	if err != nil {
		slog.Error("Could not create directory.", "log_code", "ecdb8557", "Directory", FullPath(relPath))
		return err
	}

    file, err := os.Create(FullPath(relPath, fileName))
    defer file.Close()

    if err != nil {
        slog.Error("Could not create file.", "log_code", "3b5806d7", "File path", FullPath(relPath, fileName))
    }

    return nil
}

func CreateDir(relPath string) error {
	path := FullPath(relPath)
	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		slog.Error("Could not create directory.", "log_code", "791e81d6", "Directory", path)
		return err
	}

    return nil
}

func createRootDir() error {
    return CreateDir("")
}

func GetRelDir(relPath string) (fs.FS, error) {
	path := FullPath(relPath)
    CreateDir(relPath)

	slog.Info("Created/retrieved dir.", "log_code", "17a6838b", "path", path)
	return os.DirFS(path), nil
}

func FullPath(relPath ...string) string {
    return filepath.Join(home, LimeDir, filepath.Join(relPath...))
}
