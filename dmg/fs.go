// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package dmg

import (
	"os"

	"github.com/tiagomelo/macos-dmg-creator/fs"
)

// fsOpsProvider is a variable that holds the function
// that performs file system operations.
var fsOpsProvider fsOps = defaultFsOps{}

// fsOps defines an interface for file system operations.
type fsOps interface {
	// MkdirAll creates a directory and all necessary parent directories.
	MkdirAll(path string, perm os.FileMode) error

	// CopyFile copies a file from src to dst.
	CopyFile(src, dst string) error

	// FileExists checks if a file exists at the given path.
	FileExists(path string) (bool, error)

	// CopyDir copies a directory from src to dst.
	CopyDir(src, dst string) error

	// DeleteDir deletes a directory at the given path.
	DeleteDir(path string) error

	// WriteFile writes data to a file named name.
	WriteFile(name string, data []byte, perm os.FileMode) error

	// CreateSymlink creates a symbolic link from src to dst.
	CreateSymlink(src, dst string) error
}

// defaultFsOps is the default implementation of fsOps.
type defaultFsOps struct{}

func (d defaultFsOps) MkdirAll(path string, perm os.FileMode) error {
	return fs.MkdirAll(path, perm)
}

func (d defaultFsOps) CopyFile(src, dst string) error {
	return fs.CopyFile(src, dst)
}

func (d defaultFsOps) FileExists(path string) (bool, error) {
	return fs.FileExists(path)
}

func (d defaultFsOps) CopyDir(src, dst string) error {
	return fs.CopyDir(src, dst)
}

func (d defaultFsOps) DeleteDir(path string) error {
	return fs.DeleteDir(path)
}

func (d defaultFsOps) WriteFile(name string, data []byte, perm os.FileMode) error {
	return fs.WriteFile(name, data, perm)
}

func (d defaultFsOps) CreateSymlink(src, dst string) error {
	return fs.CreateSymlink(src, dst)
}
