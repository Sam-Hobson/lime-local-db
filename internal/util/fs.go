package util

import (
	"log/slog"
	"os"
	"path/filepath"
)

type relativeFsManager struct {
	path string
}

func NewRelativeFsManager(path ...string) relativeFsManager {
	return relativeFsManager{path: filepath.Join(path...)}
}

func (relFs *relativeFsManager) FileExists(relPath string, fileName string) (bool, error) {
	_, err := os.Stat(relFs.FullPath(relPath, fileName))

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		slog.Warn("Could not check if file exists.", "log_code", "d74bd3e7", "path", relFs.FullPath(relPath, fileName))
		return false, err
	}

	return true, nil
}

func (relFs *relativeFsManager) CreateFile(relPath string, fileName string) error {
	slog.Info("Creating file.", "log_code", "12029083", "Path", relFs.FullPath(relPath, fileName))

	err := relFs.CreateDir(relPath)

	if err != nil {
		slog.Warn("Could not create directory.", "log_code", "ecdb8557", "Directory", relFs.FullPath(relPath))
		return err
	}

	file, err := os.OpenFile(relFs.FullPath(relPath, fileName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	defer file.Close()

	if err != nil {
		slog.Warn("Could not create file.", "log_code", "3b5806d7", "File path", relFs.FullPath(relPath, fileName))
		return err
	}

	slog.Info("Successfully created file.", "log_code", "31512c87", "Path", relFs.FullPath(relPath, fileName))
	return nil
}

func (relFs *relativeFsManager) MoveFile(fromRelPath, fromFileName, toRelPath, toFileName string) error {
	from := relFs.FullPath(fromRelPath, fromFileName)
	to := relFs.FullPath(toRelPath, toFileName)

	slog.Info("Moving file.", "log_code", "628970f4", "From", from, "To", to)

	if exists, err := relFs.FileExists(fromRelPath, fromFileName); err != nil {
		return err
	} else if !exists {
		slog.Warn("Cannot move file as it doesn't exist.", "log_code", "80331077", "From", from)
		return os.ErrNotExist
	}

	if err := relFs.CreateDir(toRelPath); err != nil {
		return err
	}

	if err := os.Rename(from, to); err != nil {
		slog.Warn("Could not rename file.", "log_code", "7213ca47", "From", from, "To", to)
		return err
	}

	slog.Info("Successfully moved file.", "log_code", "f8d98f03", "From", from, "To", to)
	return nil
}

func (relFs *relativeFsManager) CopyFile(fromRelPath, fromFileName, toRelPath, toFileName string) error {
	from := relFs.FullPath(fromRelPath, fromFileName)
	to := relFs.FullPath(toRelPath, toFileName)

	slog.Info("Copying file.", "log_code", "ec504e3f", "From", from, "To", to)

	if err := relFs.CreateFile(toRelPath, toFileName); err != nil {
		return err
	}

	fromData, err := os.ReadFile(from)
	if err != nil {
		slog.Warn("Could not read data from file.", "log_code", "29ab2f4a", "From", from)
		return err
	}

	err = os.WriteFile(to, fromData, os.ModePerm)
	if err != nil {
		slog.Warn("Could not write data to file.", "log_code", "2448d274", "To", to)
		return err
	}

	slog.Info("Successfully copied file.", "log_code", "ad390e7b", "From", from, "To", to)
	return nil
}

func (relFs *relativeFsManager) RmFile(relPath string, fileName string) error {
	slog.Info("Removing file.", "log_code", "017baa09", "Path", relFs.FullPath(relPath, fileName))

	err := os.Remove(relFs.FullPath(relPath, fileName))

	if err != nil {
		slog.Warn("Could not remove file.", "log_code", "8ee71301", "Path", relFs.FullPath(relPath, fileName))
		return err
	}

	slog.Info("Successfully removed file.", "log_code", "77a8feff", "Path", relFs.FullPath(relPath, fileName))
	return nil
}

func (relFs *relativeFsManager) CreateDir(relPath string) error {
	slog.Info("Creating directory.", "log_code", "109b2abd", "Path", relFs.FullPath(relPath))

	path := relFs.FullPath(relPath)
	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		slog.Warn("Could not create directory.", "log_code", "791e81d6", "Directory", path)
		return err
	}

	slog.Info("Successfully created directory.", "log_code", "88a7b6fa", "Path", relFs.FullPath(relPath))
	return nil
}

func (relFs *relativeFsManager) ReadDir(relPath string) ([]os.DirEntry, error) {
	slog.Info("Reading directory.", "log_code", "ca2b2793", "Path", relPath)

	path := relFs.FullPath(relPath)
	res, err := os.ReadDir(path)

	if err != nil {
		slog.Warn("Could not read directory.", "log_code", "76a8af95", "FullPath", path)
		return nil, err
	}

	return res, nil
}

func (relFs *relativeFsManager) FullPath(relPath ...string) string {
	return filepath.Join(relFs.path, filepath.Join(relPath...))
}
