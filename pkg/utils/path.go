package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func CheckDir(path string) error {
	stat, statErr := os.Stat(path)
	if statErr == nil {
		if !stat.IsDir() {
			return fmt.Errorf("a megadott útvonal nem egy könyvtárra mutat")
		}
	} else {
		if errors.Is(statErr, os.ErrNotExist) {
			return os.MkdirAll(path, DMode)
		}
		return statErr
	}
	return nil
}

func IsTorrentFile(path string) bool {
	if stat, statErr := os.Stat(path); errors.Is(statErr, os.ErrNotExist) {
		return false
	} else {
		if stat.IsDir() {
			return false
		}
	}
	if filepath.Ext(path) != torrentExt {
		return false
	}
	return true
}
