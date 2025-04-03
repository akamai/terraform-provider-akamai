// Package files contains utility code for working with files
package files

import (
	"os"
)

// IsSymlink returns whether the given path is a symlink.
func IsSymlink(path string) (bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	return info.Mode()&os.ModeSymlink != 0, nil
}

// IsSymlinkToDir returns whether the given path is a symlink pointing to a directory.
func IsSymlinkToDir(path string) (bool, error) {
	isLink, err := IsSymlink(path)
	if err != nil {
		return false, err
	}
	if !isLink {
		return false, nil
	}
	// os.Stat follows a symlink when given one and returns information about the target file
	targetInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return targetInfo.IsDir(), nil
}
