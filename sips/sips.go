// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package sips

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/tiagomelo/macos-dmg-creator/syscall"
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

// GenerateIcons generates icons with the specified sizes.
func GenerateIcons(iconPath, outputDir string, sizes ...int) error {
	for _, size := range sizes {
		if _, err := osCommandExecutorProvider.ExecCommand("sips", "-z", fmt.Sprintf("%d", size), fmt.Sprintf("%d", size), iconPath, "--out", fmt.Sprintf("%s/icon_%dx%d.png", outputDir, size, size)); err != nil {
			return errors.Wrapf(err, "error when generating icon with size %d", size)
		}
	}
	return nil
}
