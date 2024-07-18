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

var home = os.Getenv("HOME")

func FileExists(relPath string, fileName string) (bool, error) {
    _, err := os.Stat(FullPath(relPath, fileName))

    if err != nil {
        if os.IsNotExist(err) {
            return false, nil
        }

        slog.Error("Could not check if file exists.", "log_code", "d74bd3e7", "path", FullPath(relPath, fileName))
        return false, err
    }

    return true, nil
}

func CreateFile(relPath string, fileName string) error {
    slog.Info("Creating file.", "log_code", "12029083", "Path", FullPath(relPath, fileName))

	err := CreateDir(relPath)

	if err != nil {
		slog.Error("Could not create directory.", "log_code", "ecdb8557", "Directory", FullPath(relPath))
		return err
	}

    file, err := os.OpenFile(FullPath(relPath, fileName), os.O_RDWR | os.O_CREATE | os.O_TRUNC, os.ModePerm)
    defer file.Close()

    if err != nil {
        slog.Error("Could not create file.", "log_code", "3b5806d7", "File path", FullPath(relPath, fileName))
    }

    slog.Info("Successfully created file.", "log_code", "31512c87", "Path", FullPath(relPath, fileName))
    return nil
}

func MoveFile(fromRelPath, fromFileName, toRelPath, toFileName string) error {
    from := FullPath(fromRelPath, fromFileName)
    to := FullPath(toRelPath, toFileName)

    slog.Info("Moving file.", "log_code", "628970f4", "From", from, "To", to)
    exists, err := FileExists(fromRelPath, fromFileName)

    if err != nil {
        return err
    }

    if !exists {
        slog.Error("Cannot move file as it doesn't exist.", "log_code", "80331077", "From", from)
        return os.ErrNotExist
    }

    if err := CreateDir(toRelPath); err != nil {
        return err
    }

    if err := os.Rename(from, to); err != nil {
        slog.Error("Could not rename file.", "log_code", "7213ca47", "From", from, "To", to)
        return err
    }

    slog.Info("Successfully moved file.", "log_code", "f8d98f03", "From", from, "To", to)
    return nil
}

func RmFile(relPath string, fileName string) error {
    slog.Info("Removing file.", "log_code", "017baa09", "Path", FullPath(relPath, fileName))

    err := os.Remove(FullPath(relPath, fileName))

    if err != nil {
        slog.Error("Could not remove file.", "log_code", "8ee71301", "Path", FullPath(relPath, fileName))
        return err
    }

    slog.Info("Successfully removed file.", "log_code", "77a8feff", "Path", FullPath(relPath, fileName))
    return nil
}

func CreateDir(relPath string) error {
    slog.Info("Creating directory.", "log_code", "109b2abd", "Path", FullPath(relPath))

	path := FullPath(relPath)
	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		slog.Error("Could not create directory.", "log_code", "791e81d6", "Directory", path)
		return err
	}

    slog.Info("Successfully created directory.", "log_code", "88a7b6fa", "Path", FullPath(relPath))
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
