/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package linux

import (
	"os"
	"os/user"
	"path/filepath"
	"syscall"
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	R_OK = 4
	W_OK = 2
	X_OK = 1
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Returns path to user home directory
func UserDir() string {
	currentUser, _ := user.Current()
	return currentUser.HomeDir
}

// Make absolute path from a path, relative to another
// Returns the absolute path and a boolean value which indicates
// if the returned path exists or not
func ResolvePath(path string, relpath string) (string, bool) {

	// Deal with ~/ form - substitute user's home path
	if filepath.HasPrefix(path, "~/") {
		path = filepath.Join(UserDir(), path[2:])
	}

	// Join relpath with path
	if filepath.IsAbs(path) == false {
		path = filepath.Join(relpath, path)
	}

	// Clean up the path
	path = filepath.Clean(path)

	// Determine if path exists, return
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path, false
	} else {
		return path, true
	}
}

// IsWritableFolder checks a folder for writable
func IsWritableFolder(path string) bool {
	if stat, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else if stat.IsDir() == false {
		return false
	} else {
		return syscall.Access(path, W_OK) == nil
	}
}
