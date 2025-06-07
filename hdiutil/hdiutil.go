// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package hdiutil

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/tiagomelo/go-retry"
	"github.com/tiagomelo/macos-dmg-creator/syscall"
)

// osStat is a variable that holds the function
// that retrieves file information.
var osStat = os.Stat

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

// CreateDMG creates a dmg file.
func CreateDMG(size, fs, volName, layout, output string) error {
	if _, err := osCommandExecutorProvider.ExecCommand("hdiutil", "create", "-size",
		size, "-fs", fs, "-volname", volName, "-layout",
		layout, "-o", output); err != nil {
		return errors.Wrapf(err, "error when creating dmg with name %s", volName)
	}
	return nil
}

// MountDMG mounts a dmg file.
// The dmg image is mounted at /Volumes/volName.
func MountDMG(dmgVolName, dmgPath string) error {
	if _, err := osCommandExecutorProvider.ExecCommand("hdiutil", "attach", dmgPath); err != nil {
		return errors.Wrapf(err, "error when attaching dmg with name %s", dmgPath)
	}

	// check if the volume is mounted
	// and retry if it is not mounted yet.
	dmgVolName = fmt.Sprintf("/Volumes/%s", dmgVolName)
	linearBackoffStrategy := retry.NewLinearBackoff(100*time.Millisecond, 1*time.Second, 10)
	_, err := retry.Do(func() error {
		if _, err := osStat(dmgVolName); err != nil {
			return err
		}
		return nil
	}, linearBackoffStrategy)
	if err != nil {
		return errors.Wrapf(err, "error when waiting for volume %s to be mounted", dmgPath)
	}

	return nil
}

// UnmountDMG unmounts a dmg file
// that was previously attached to /Volumes/volName.
func UnmountDMG(volName string) error {
	if _, err := osCommandExecutorProvider.ExecCommand("hdiutil", "detach", volName); err != nil {
		return errors.Wrapf(err, "error when detaching dmg with name %s", volName)
	}

	// check if the volume is unmounted
	// and retry if it is still present.
	strategy := retry.NewLinearBackoff(100*time.Millisecond, 1*time.Second, 10)
	_, err := retry.Do(func() error {
		_, err := osStat(volName)
		if os.IsNotExist(err) {
			return nil // success, volume is gone.
		}
		return err
	}, strategy)
	if err != nil {
		return errors.Wrapf(err, "error when waiting for volume %s to be unmounted", volName)
	}

	return nil
}

// ConvertDMG converts a dmg file.
// It aims to convert it to a compressed format.
func ConvertDMG(dmgPath, dmgOutputFileName string) error {
	if _, err := osCommandExecutorProvider.ExecCommand("hdiutil", "convert", dmgPath, "-format", "UDZO", "-o", dmgOutputFileName); err != nil {
		return errors.Wrapf(err, "error when converting dmg file %s", dmgPath)
	}
	return nil
}
