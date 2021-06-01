package ioutils

import (
	"os"
	"path"
	"strings"
)

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func PureName(name string) string {
	return strings.TrimSuffix(name, path.Ext(name))
}

func IsDir(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			return true
		}
	} else {
		return info.IsDir()
	}
}
