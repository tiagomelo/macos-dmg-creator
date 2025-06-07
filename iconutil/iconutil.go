// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package iconutil

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

// GenerateIconSet generates an icon set from the specified icons directory.
func GenerateIconSet(iconsDir, outputDir string) error {
	outputDir = fmt.Sprintf("%s/icon.icns", outputDir)
	if _, err := osCommandExecutorProvider.ExecCommand("iconutil", "-c", "icns", "-o", outputDir, iconsDir); err != nil {
		return errors.Wrap(err, "error when generating icon set")
	}
	return nil
}
