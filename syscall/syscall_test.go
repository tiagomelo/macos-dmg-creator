// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package syscall

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExecCommand(t *testing.T) {
	testCases := []struct {
		name            string
		mockExecCommand func(name string, arg ...string) ([]byte, error)
		expectedOutput  []byte
		expectedError   error
	}{
		{
			name: "happy path",
			mockExecCommand: func(name string, arg ...string) ([]byte, error) {
				return []byte("output"), nil
			},
			expectedOutput: []byte("output"),
		},
		{
			name: "error",
			mockExecCommand: func(name string, arg ...string) ([]byte, error) {
				return []byte("output"), errors.New("some error")
			},
			expectedError: errors.New("error when executing command [ls] with args [-l]: output: [output]: some error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			execCommand = tc.mockExecCommand
			output, err := ExecCommand("ls", "-l")
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf(`expected error "%v", got nil`, tc.expectedError)
				}
				require.Equal(t, string(tc.expectedOutput), string(output))
			}
		})
	}
}
