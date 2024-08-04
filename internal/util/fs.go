package util

import (
	"io"
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

		slog.Warn("Could not check if file exists.", "Log code", "d74bd3e7", "Path", relFs.FullPath(relPath, fileName))
		return false, err
	}

	return true, nil
}

func (relFs *relativeFsManager) CreateFile(relPath string, fileName string) error {
	slog.Info("Creating file.", "Log code", "12029083", "Path", relFs.FullPath(relPath, fileName))

	if err := relFs.CreateDir(relPath); err != nil {
		slog.Warn("Could not create directory.", "Log code", "ecdb8557", "Directory", relFs.FullPath(relPath))
		return err
	}

	if file, err := os.OpenFile(relFs.FullPath(relPath, fileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm); err != nil {
		slog.Warn("Could not create file.", "Log code", "3b5806d7", "File path", relFs.FullPath(relPath, fileName))
		return err
	} else {
		file.Close()
	}

	slog.Info("Successfully created file.", "Log code", "31512c87", "Path", relFs.FullPath(relPath, fileName))
	return nil
}

func (relFs *relativeFsManager) OpenFile(relPath string, fileName string) (io.Writer, error) {
	slog.Info("Opening file.", "Log code", "6453d64c", "Path", relFs.FullPath(relPath, fileName))

	if err := relFs.CreateFile(relPath, fileName); err != nil {
		return nil, err
	}

	return os.OpenFile(relFs.FullPath(relPath, fileName), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
}

func (relFs *relativeFsManager) MoveFile(fromRelPath, fromFileName, toRelPath, toFileName string) error {
	from := relFs.FullPath(fromRelPath, fromFileName)
	to := relFs.FullPath(toRelPath, toFileName)

	slog.Info("Moving file.", "Log code", "628970f4", "From", from, "To", to)

	if exists, err := relFs.FileExists(fromRelPath, fromFileName); err != nil {
		return err
	} else if !exists {
		slog.Warn("Cannot move file as it doesn't exist.", "Log code", "80331077", "From", from)
		return os.ErrNotExist
	}

	if err := relFs.CreateDir(toRelPath); err != nil {
		return err
	}

	if err := os.Rename(from, to); err != nil {
		slog.Warn("Could not rename file.", "Log code", "7213ca47", "From", from, "To", to)
		return err
	}

	slog.Info("Successfully moved file.", "Log code", "f8d98f03", "From", from, "To", to)
	return nil
}

func (relFs *relativeFsManager) CopyFile(fromRelPath, fromFileName, toRelPath, toFileName string) error {
	from := relFs.FullPath(fromRelPath, fromFileName)
	to := relFs.FullPath(toRelPath, toFileName)

	slog.Info("Copying file.", "Log code", "ec504e3f", "From", from, "To", to)

	if err := relFs.CreateFile(toRelPath, toFileName); err != nil {
		return err
	}

	if fromData, err := os.ReadFile(from); err != nil {
		slog.Warn("Could not read data from file.", "Log code", "29ab2f4a", "From", from)
		return err
	} else if err := os.WriteFile(to, fromData, os.ModePerm); err != nil {
		slog.Warn("Could not write data to file.", "Log code", "2448d274", "To", to)
		return err
	}

	slog.Info("Successfully copied file.", "Log code", "ad390e7b", "From", from, "To", to)
	return nil
}

func (relFs *relativeFsManager) RmFile(relPath string, fileName string) error {
	slog.Info("Removing file.", "Log code", "017baa09", "Path", relFs.FullPath(relPath, fileName))

	err := os.Remove(relFs.FullPath(relPath, fileName))

	if err != nil {
		slog.Warn("Could not remove file.", "Log code", "8ee71301", "Path", relFs.FullPath(relPath, fileName))
		return err
	}

	slog.Info("Successfully removed file.", "Log code", "77a8feff", "Path", relFs.FullPath(relPath, fileName))
	return nil
}

func (relFs *relativeFsManager) CreateDir(relPath string) error {
	slog.Info("Creating directory.", "Log code", "109b2abd", "Path", relFs.FullPath(relPath))

	path := relFs.FullPath(relPath)
	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		slog.Warn("Could not create directory.", "Log code", "791e81d6", "Directory", path)
		return err
	}

	slog.Info("Successfully created directory.", "Log code", "88a7b6fa", "Path", relFs.FullPath(relPath))
	return nil
}

func (relFs *relativeFsManager) ReadDir(relPath string) ([]os.DirEntry, error) {
	slog.Info("Reading directory.", "Log code", "ca2b2793", "Path", relPath)

	path := relFs.FullPath(relPath)
	res, err := os.ReadDir(path)

	if err != nil {
		slog.Warn("Could not read directory.", "Log code", "76a8af95", "Full path", path)
		return nil, err
	}

	return res, nil
}

func (relFs *relativeFsManager) FullPath(relPath ...string) string {
	return filepath.Join(relFs.path, filepath.Join(relPath...))
}
