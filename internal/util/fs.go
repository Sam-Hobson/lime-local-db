package util

import (
	"os"
	"path/filepath"
)

type relativeFsManager struct {
	path string
}

func NewRelativeFsManager(path ...string) relativeFsManager {
	if path == nil {
		return relativeFsManager{path: ""}
	}
	return relativeFsManager{path: filepath.Join(path...)}
}

func (relFs *relativeFsManager) FileExists(relFilepath ...string) (bool, error) {
	fullpath := relFs.FullPath(relFilepath...)
	_, err := os.Stat(fullpath)

	Log("255b1802").Info("Checking if file exists.", "File", fullpath)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		Log("d74bd3e7").Warn("Could not check if file exists.", "Path", fullpath)
		return false, err
	}

	return true, nil
}

func (relFs *relativeFsManager) CreateFile(relFilepath ...string) error {
	fullpath := relFs.FullPath(relFilepath...)
	Log("0c2f42ae").Info("Creating file.", "File", fullpath)

	if err := relFs.CreateDir(relFs.Dir(relFilepath...)); err != nil {
		return err
	}

	if file, err := os.OpenFile(fullpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm); err != nil {
		Log("3b5806d7").Warn("Could not create file.", "File path", fullpath)
		return err
	} else {
		file.Close()
	}

	return nil
}

func (relFs *relativeFsManager) OpenFile(relFilepath ...string) (*os.File, error) {
	fullpath := relFs.FullPath(relFilepath...)
	Log("6453d64c").Info("Opening file.", "Path", fullpath)

	if err := relFs.CreateFile(relFilepath...); err != nil {
		return nil, err
	}

	return os.OpenFile(fullpath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
}

func (relFs *relativeFsManager) MoveFile(relFrompath, relTopath string) error {
	from := relFs.FullPath(relFrompath)
	to := relFs.FullPath(relTopath)

	Log("628970f4").Info("Moving file.", "From", from, "To", to)

	if exists, err := relFs.FileExists(relFrompath); err != nil {
		return err
	} else if !exists {
		Log("80331077").Warn("Cannot move file as it doesn't exist.", "From", from)
		return os.ErrNotExist
	}

	if err := relFs.CreateDir(relFs.Dir(relTopath)); err != nil {
		return err
	}

	if err := os.Rename(from, to); err != nil {
		Log("7213ca47").Warn("Could not rename file.", "From", from, "To", to)
		return err
	}

	return nil
}

func (relFs *relativeFsManager) CopyFile(relFrompath, relTopath string) error {
	from := relFs.FullPath(relFrompath)
	to := relFs.FullPath(relTopath)

	Log("ec504e3f").Info("Copying file.", "From", from, "To", to)

	if err := relFs.CreateFile(relTopath); err != nil {
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

func (relFs *relativeFsManager) RmFile(relFilepath ...string) error {
	fullpath := relFs.FullPath(relFilepath...)
	Log("017baa09").Info("Removing file.", "Path", fullpath)

	err := os.Remove(fullpath)

	if err != nil {
		Log("8ee71301").Warn("Could not remove file.", "Path", fullpath)
		return err
	}

	return nil
}

func (relFs *relativeFsManager) CreateDir(relPath ...string) error {
	fullpath := relFs.FullPath(relPath...)
	Log("109b2abd").Info("Creating directory.", "Path", fullpath)

	err := os.MkdirAll(fullpath, os.ModePerm)

	if err != nil {
		Log("791e81d6").Warn("Could not create directory.", "Directory", fullpath)
		return err
	}

	return nil
}

func (relFs *relativeFsManager) ReadDir(relPath ...string) ([]os.DirEntry, error) {
	fullpath := relFs.FullPath(relPath...)
	Log("ca2b2793").Info("Reading directory.", "Path", fullpath)

	res, err := os.ReadDir(fullpath)

	if err != nil {
		Log("76a8af95").Warn("Could not read directory.", "Full path", fullpath)
		return nil, err
	}

	return res, nil
}

func (relFs *relativeFsManager) ReadFileIntoMemry(relFilepath ...string) (string, error) {
	fullpath := relFs.FullPath(relFilepath...)
	Log("2133dd24").Info("Reading file.", "Path", fullpath)

	data, err := os.ReadFile(fullpath)
	if err != nil {
		Log("cc1e0022").Warn("Could not read file into memory.", "Path", fullpath)
		return "", err
	}

	return string(data), nil
}

func (relFs *relativeFsManager) Dir(relPath ...string) string {
	return filepath.Dir(filepath.Join(relPath...))
}

func (relFs *relativeFsManager) FullDir(relPath ...string) string {
	return filepath.Dir(relFs.FullPath(relPath...))
}

func (relFs *relativeFsManager) FullPath(relPath ...string) string {
	return filepath.Join(relFs.path, filepath.Join(relPath...))
}
