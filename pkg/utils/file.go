package utils

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func FileExist(basePath, filePath string) (bool, error) {
	fullPath, err := GetFullPath(basePath, filePath)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(fullPath)
	return err == nil, nil
}

func GetFullPath(basePath, filePath string) (string, error) {
	var err error
	if basePath == "" {
		basePath, err = os.Getwd()
		if err != nil {
			return "", errors.Wrap(err, "Error gett pwd path")
		}
	}
	fullPath, err := filepath.Abs(filepath.Join(basePath, filePath))
	if err != nil {
		return "", errors.Wrap(err, "Error abs path getting")
	}
	return fullPath, err
}
