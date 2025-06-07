// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package syscall

import (
	"os/exec"

	"github.com/pkg/errors"
)

// execCommand is a variable that holds the function that executes a command with arguments.
// It is a variable so that it can be mocked in tests.
var execCommand = func(name string, arg ...string) ([]byte, error) {
	return exec.Command(name, arg...).CombinedOutput()
}

// ExecCommand executes a command with arguments.
func ExecCommand(cmd string, args ...string) (string, error) {
	output, err := execCommand(cmd, args...)
	if err != nil {
		return "", errors.Wrapf(err, "error when executing command [%s] with args %v: output: [%v]", cmd, args, string(output))
	}
	return string(output), nil
}
