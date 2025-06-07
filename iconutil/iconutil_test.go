// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package iconutil

import (
	"testing"

	"github.com/pkg/errors"
)

func TestGenerateIconSet(t *testing.T) {
	testCases := []struct {
		name                  string
		mockOsCommandExecutor func() *mockOsCommandExecutor
		expectedError         error
	}{
		{
			name: "happy path",
			mockOsCommandExecutor: func() *mockOsCommandExecutor {
				return &mockOsCommandExecutor{}
			},
		},
		{
			name: "error",
			mockOsCommandExecutor: func() *mockOsCommandExecutor {
				return &mockOsCommandExecutor{
					err: errors.New("some error"),
				}
			},
			expectedError: errors.New("error when generating icon set: some error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			osCommandExecutorProvider = tc.mockOsCommandExecutor()
			err := GenerateIconSet("iconsDir", "outputDir")
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				if tc.expectedError.Error() != err.Error() {
					t.Fatalf(`expected error "%v", got "%v"`, tc.expectedError, err)
				}
			} else {
				if tc.expectedError != nil {
					t.Fatalf(`expected error "%v", got nil`, tc.expectedError)
				}
			}
		})
	}
}

type mockOsCommandExecutor struct {
	err error
}

func (m *mockOsCommandExecutor) ExecCommand(name string, arg ...string) (string, error) {
	return "", m.err
}
