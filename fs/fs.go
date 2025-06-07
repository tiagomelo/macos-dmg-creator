// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package fs

import (
	"os"

	"github.com/pkg/errors"
	"github.com/tiagomelo/macos-dmg-creator/syscall"
)

// for ease of unit testing.
var (
	osStat      = os.Stat
	osWriteFile = os.WriteFile
)

// osCommandExecutorProvider is a variable that holds the function
// that executes a command with arguments.
var osCommandExecutorProvider osCommandExecutor = &defaultOsCommandExecutor{}

// osCommandExecutor defines an interface for executing OS commands.
type osCommandExecutor interface {
	ExecCommand(name string, arg ...string) (string, error)
}

// defaultOsCommandExecutor is the default implementation of osCommandExecutor.
type defaultOsCommandExecutor struct{}

// ExecCommand executes a command with arguments.
func (d *defaultOsCommandExecutor) ExecCommand(name string, arg ...string) (string, error) {
	return syscall.ExecCommand(name, arg...)
}

// CopyFile copies a file from src to dst.
func CopyFile(src, dst string) error {
	if _, err := osCommandExecutorProvider.ExecCommand("cp", src, dst); err != nil {
		return errors.Wrapf(err, "error when copying file from [%s] to [%s]", src, dst)
	}
	return nil
}

// CopyDir copies a directory from src to dst.
func CopyDir(src, dst string) error {
	if _, err := osCommandExecutorProvider.ExecCommand("cp", "-r", src, dst); err != nil {
		return errors.Wrapf(err, "error when copying directory from [%s] to [%s]", src, dst)
	}
	return nil
}

// FileExists checks if a file exists.
func FileExists(path string) (bool, error) {
	info, err := osStat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrapf(err, "error when checking if [%s] exists", path)
	}
	if info.IsDir() {
		return false, errors.Errorf("[%s] is a directory, not a file", path)
	}
	return true, nil
}

// DeleteFile deletes a file.
func DeleteFile(path string) error {
	if _, err := osCommandExecutorProvider.ExecCommand("rm", "-f", path); err != nil {
		return errors.Wrapf(err, "error when deleting file [%s]", path)
	}
	return nil
}

// DeleteDir deletes a directory.
func DeleteDir(path string) error {
	if _, err := osCommandExecutorProvider.ExecCommand("rm", "-rf", path); err != nil {
		return errors.Wrapf(err, "error when deleting directory [%s]", path)
	}
	return nil
}

// MkdirAll creates a directory and all necessary parents.
func MkdirAll(path string, perm os.FileMode) error {
	if _, err := osCommandExecutorProvider.ExecCommand("mkdir", "-p", path); err != nil {
		return errors.Wrapf(err, "error when creating directory path [%s]", path)
	}
	return nil
}

// WriteFile writes data to a file named name.
func WriteFile(name string, data []byte, perm os.FileMode) error {
	if err := osWriteFile(name, data, perm); err != nil {
		return errors.Wrapf(err, "error when writing to file [%s]", name)
	}
	return nil
}

// CreateSymlink creates a symbolic link from src to dst.
func CreateSymlink(src, dst string) error {
	if _, err := osCommandExecutorProvider.ExecCommand("ln", "-s", src, dst); err != nil {
		return errors.Wrapf(err, "error when creating symlink from [%s] to [%s]", src, dst)
	}
	return nil
}
