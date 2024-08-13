package util

import (
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

		Log("d74bd3e7").Warn("Could not check if file exists.", "Path", relFs.FullPath(relPath, fileName))
		return false, err
	}

	return true, nil
}

func (relFs *relativeFsManager) CreateFile(relPath string, fileName string) error {
	Log("0c2f42ae").Info("Creating file.", "Path", relFs.FullPath(relPath, fileName))

	if err := relFs.CreateDir(relPath); err != nil {
		Log("ecdb8557").Warn("Could not create directory.", "Directory", relFs.FullPath(relPath))
		return err
	}

	if file, err := os.OpenFile(relFs.FullPath(relPath, fileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm); err != nil {
		Log("3b5806d7").Warn("Could not create file.", "File path", relFs.FullPath(relPath, fileName))
		return err
	} else {
		file.Close()
	}

	return nil
}

func (relFs *relativeFsManager) OpenFile(relPath string, fileName string) (*os.File, error) {
	Log("6453d64c").Info("Opening file.", "Path", relFs.FullPath(relPath, fileName))

	if err := relFs.CreateFile(relPath, fileName); err != nil {
		return nil, err
	}

	return os.OpenFile(relFs.FullPath(relPath, fileName), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
}

func (relFs *relativeFsManager) MoveFile(fromRelPath, fromFileName, toRelPath, toFileName string) error {
	from := relFs.FullPath(fromRelPath, fromFileName)
	to := relFs.FullPath(toRelPath, toFileName)

	Log("628970f4").Info("Moving file.", "From", from, "To", to)

	if exists, err := relFs.FileExists(fromRelPath, fromFileName); err != nil {
		return err
	} else if !exists {
		Log("80331077").Warn("Cannot move file as it doesn't exist.", "From", from)
		return os.ErrNotExist
	}

	if err := relFs.CreateDir(toRelPath); err != nil {
		return err
	}

	if err := os.Rename(from, to); err != nil {
		Log("7213ca47").Warn("Could not rename file.", "From", from, "To", to)
		return err
	}

	return nil
}

func (relFs *relativeFsManager) CopyFile(fromRelPath, fromFileName, toRelPath, toFileName string) error {
	from := relFs.FullPath(fromRelPath, fromFileName)
	to := relFs.FullPath(toRelPath, toFileName)

	Log("ec504e3f").Info("Copying file.", "From", from, "To", to)

	if err := relFs.CreateFile(toRelPath, toFileName); err != nil {
		return err
	}

	if fromData, err := os.ReadFile(from); err != nil {
		Log("29ab2f4a").Warn("Could not read data from file.", "From", from)
		return err
	} else if err := os.WriteFile(to, fromData, os.ModePerm); err != nil {
		Log("2448d274").Warn("Could not write data to file.", "To", to)
		return err
	}

	return nil
}

func (relFs *relativeFsManager) RmFile(relPath string, fileName string) error {
	Log("017baa09").Info("Removing file.", "Path", relFs.FullPath(relPath, fileName))

	err := os.Remove(relFs.FullPath(relPath, fileName))

	if err != nil {
		Log("8ee71301").Warn("Could not remove file.", "Path", relFs.FullPath(relPath, fileName))
		return err
	}

	return nil
}

func (relFs *relativeFsManager) CreateDir(relPath string) error {
	Log("109b2abd").Info("Creating directory.", "Path", relFs.FullPath(relPath))

	path := relFs.FullPath(relPath)
	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		Log("791e81d6").Warn("Could not create directory.", "Directory", path)
		return err
	}

	return nil
}

func (relFs *relativeFsManager) ReadDir(relPath string) ([]os.DirEntry, error) {
	Log("ca2b2793").Info("Reading directory.", "Path", relPath)

	path := relFs.FullPath(relPath)
	res, err := os.ReadDir(path)

	if err != nil {
		Log("76a8af95").Warn("Could not read directory.", "Full path", path)
		return nil, err
	}

	return res, nil
}

func (relFs *relativeFsManager) ReadFileIntoMemry(relPath string, fileName string) (string, error) {
	path := relFs.FullPath(relPath, fileName)
	Log("2133dd24").Info("Reading file.", "Path", path)

	data, err := os.ReadFile(path)
	if err != nil {
		Log("cc1e0022").Warn("Could not read file into memory.", "Path", path)
		return "", err
	}

	return string(data), nil
}

func (relFs *relativeFsManager) FullPath(relPath ...string) string {
	return filepath.Join(relFs.path, filepath.Join(relPath...))
}
