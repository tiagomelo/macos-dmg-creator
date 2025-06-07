// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package hdiutil

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateDmg(t *testing.T) {
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
			expectedError: errors.New("error when creating dmg with name Test: some error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			osCommandExecutorProvider = tc.mockOsCommandExecutor()
			err := CreateDMG("1m", "HFS+", "Test", "Standard", "test.dmg")
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf(`expected error "%v", got nil`, tc.expectedError)
				}
			}
		})
	}
}

func TestMountDMG(t *testing.T) {
	testCases := []struct {
		name                  string
		mockOsCommandExecutor func() *mockOsCommandExecutor
		mockOsStat            func(name string) (os.FileInfo, error)
		expectedError         error
	}{
		{
			name: "happy path",
			mockOsCommandExecutor: func() *mockOsCommandExecutor {
				return &mockOsCommandExecutor{}
			},
			mockOsStat: func(name string) (os.FileInfo, error) {
				return &mockFileInfo{name: "Test"}, nil
			},
		},
		{
			name: "fail to mount dmg",
			mockOsCommandExecutor: func() *mockOsCommandExecutor {
				return &mockOsCommandExecutor{
					err: errors.New("some error"),
				}
			},
			mockOsStat: func(name string) (os.FileInfo, error) {
				return &mockFileInfo{}, nil
			},
			expectedError: errors.New("error when attaching dmg with name Test: some error"),
		},
		{
			name: "fail to check mounted volume",
			mockOsCommandExecutor: func() *mockOsCommandExecutor {
				return &mockOsCommandExecutor{}
			},
			mockOsStat: func(name string) (os.FileInfo, error) {
				return nil, errors.New("some error")
			},
			expectedError: errors.New("error when waiting for volume Test to be mounted: some error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			osStat = tc.mockOsStat
			osCommandExecutorProvider = tc.mockOsCommandExecutor()
			err := MountDMG("", "Test")
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf(`expected error "%v", got nil`, tc.expectedError)
				}
			}
		})
	}
}

func TestUnmountDmg(t *testing.T) {
	testCases := []struct {
		name                  string
		mockOsCommandExecutor func() *mockOsCommandExecutor
		mockOsStat            func(name string) (os.FileInfo, error)
		expectedError         error
	}{
		{
			name: "happy path",
			mockOsCommandExecutor: func() *mockOsCommandExecutor {
				return &mockOsCommandExecutor{}
			},
			mockOsStat: func(name string) (os.FileInfo, error) {
				return nil, os.ErrNotExist
			},
		},
		{
			name: "fail to unmount dmg",
			mockOsCommandExecutor: func() *mockOsCommandExecutor {
				return &mockOsCommandExecutor{
					err: errors.New("some error"),
				}
			},
			expectedError: errors.New("error when detaching dmg with name Test: some error"),
		},
		{
			name: "fail to check unmounted volume",
			mockOsCommandExecutor: func() *mockOsCommandExecutor {
				return &mockOsCommandExecutor{}
			},
			mockOsStat: func(name string) (os.FileInfo, error) {
				return nil, errors.New("some error")
			},
			expectedError: errors.New("error when waiting for volume Test to be unmounted: some error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			osStat = tc.mockOsStat
			osCommandExecutorProvider = tc.mockOsCommandExecutor()
			err := UnmountDMG("Test")
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf(`expected error "%v", got nil`, tc.expectedError)
				}
			}
		})
	}
}

func TestConvertDmg(t *testing.T) {
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
			expectedError: errors.New("error when converting dmg file Test.dmg: some error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			osCommandExecutorProvider = tc.mockOsCommandExecutor()
			err := ConvertDMG("Test.dmg", "Test-converted.dmg")
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
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

type mockFileInfo struct {
	name string
}

func (f *mockFileInfo) Name() string       { return f.name }
func (f *mockFileInfo) Size() int64        { return 0 }
func (f *mockFileInfo) Mode() os.FileMode  { return 0 }
func (f *mockFileInfo) ModTime() time.Time { return time.Now() }
func (f *mockFileInfo) IsDir() bool        { return false }
func (f *mockFileInfo) Sys() interface{}   { return nil }
