// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package fs

import (
	sysFs "io/fs"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestCopyFile(t *testing.T) {
	testCases := []struct {
		name                  string
		mockOsCommandExecutor func() *mockOsCommandExecutor
		wantErr               error
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
			wantErr: errors.New("error when copying file from [src] to [dst]: some error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			osCommandExecutorProvider = tc.mockOsCommandExecutor()
			err := CopyFile("src", "dst")
			if err != nil {
				if tc.wantErr == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				if tc.wantErr.Error() != err.Error() {
					t.Fatalf(`expected error "%v", got "%v"`, tc.wantErr, err)
				}
			} else {
				if tc.wantErr != nil {
					t.Fatalf(`expected error "%v", got nil`, tc.wantErr)
				}
			}
		})
	}
}

func TestCopyDir(t *testing.T) {
	testCases := []struct {
		name                  string
		mockOsCommandExecutor func() *mockOsCommandExecutor
		wantErr               error
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
			wantErr: errors.New("error when copying directory from [src] to [dst]: some error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			osCommandExecutorProvider = tc.mockOsCommandExecutor()
			err := CopyDir("src", "dst")
			if err != nil {
				if tc.wantErr == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				if tc.wantErr.Error() != err.Error() {
					t.Fatalf(`expected error "%v", got "%v"`, tc.wantErr, err)
				}
			} else {
				if tc.wantErr != nil {
					t.Fatalf(`expected error "%v", got nil`, tc.wantErr)
				}
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	testCases := []struct {
		name       string
		mockOsStat func(name string) (sysFs.FileInfo, error)
		want       bool
		wanterror  error
	}{
		{
			name: "exists",
			mockOsStat: func(name string) (sysFs.FileInfo, error) {
				return &mockFileInfo{isDir: false}, nil
			},
			want: true,
		},
		{
			name: "does not exist",
			mockOsStat: func(name string) (sysFs.FileInfo, error) {
				return nil, os.ErrNotExist
			},
		},
		{
			name: "error",
			mockOsStat: func(name string) (sysFs.FileInfo, error) {
				return nil, errors.New("some error")
			},
			wanterror: errors.New("error when checking if [someFile] exists: some error"),
		},
		{
			name: "is a directory",
			mockOsStat: func(name string) (sysFs.FileInfo, error) {
				return &mockFileInfo{isDir: true}, nil
			},
			wanterror: errors.New("[someFile] is a directory, not a file"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			osStat = tc.mockOsStat

			output, err := FileExists("someFile")
			if err != nil {
				if tc.wanterror == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				if tc.wanterror.Error() != err.Error() {
					t.Fatalf(`expected error "%v", got "%v"`, tc.wanterror, err)
				}
			} else {
				if tc.wanterror != nil {
					t.Fatalf(`expected error "%v", got nil`, tc.wanterror)
				}
				require.Equal(t, tc.want, output)
			}
		})
	}
}

func TestDeleteFile(t *testing.T) {
	testCases := []struct {
		name                  string
		mockOsCommandExecutor func() *mockOsCommandExecutor
		wantErr               error
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
			wantErr: errors.New("error when deleting file [someFile]: some error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			osCommandExecutorProvider = tc.mockOsCommandExecutor()
			err := DeleteFile("someFile")
			if err != nil {
				if tc.wantErr == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				if tc.wantErr.Error() != err.Error() {
					t.Fatalf(`expected error "%v", got "%v"`, tc.wantErr, err)
				}
			} else {
				if tc.wantErr != nil {
					t.Fatalf(`expected error "%v", got nil`, tc.wantErr)
				}
			}
		})
	}
}

func TestDeleteDir(t *testing.T) {
	testCases := []struct {
		name                  string
		mockOsCommandExecutor func() *mockOsCommandExecutor
		wantErr               error
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
			wantErr: errors.New("error when deleting directory [someDir]: some error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			osCommandExecutorProvider = tc.mockOsCommandExecutor()
			err := DeleteDir("someDir")
			if err != nil {
				if tc.wantErr == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				if tc.wantErr.Error() != err.Error() {
					t.Fatalf(`expected error "%v", got "%v"`, tc.wantErr, err)
				}
			} else {
				if tc.wantErr != nil {
					t.Fatalf(`expected error "%v", got nil`, tc.wantErr)
				}
			}
		})
	}
}

func TestMkdirAll(t *testing.T) {
	testCases := []struct {
		name                  string
		mockOsCommandExecutor func() *mockOsCommandExecutor
		wantErr               error
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
			wantErr: errors.New("error when creating directory path [path/to/someDir]: some error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			osCommandExecutorProvider = tc.mockOsCommandExecutor()
			err := MkdirAll("path/to/someDir", os.ModePerm)
			if err != nil {
				if tc.wantErr == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				if tc.wantErr.Error() != err.Error() {
					t.Fatalf(`expected error "%v", got "%v"`, tc.wantErr, err)
				}
			} else {
				if tc.wantErr != nil {
					t.Fatalf(`expected error "%v", got nil`, tc.wantErr)
				}
			}
		})
	}
}

func TestWriteFile(t *testing.T) {
	testCases := []struct {
		name            string
		mockOsWriteFile func(name string, data []byte, perm os.FileMode) error
		wantErr         error
	}{
		{
			name: "happy path",
			mockOsWriteFile: func(name string, data []byte, perm os.FileMode) error {
				return nil
			},
		},
		{
			name: "error",
			mockOsWriteFile: func(name string, data []byte, perm os.FileMode) error {
				return errors.New("some error")
			},
			wantErr: errors.New("error when writing to file [someFile]: some error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			osWriteFile = tc.mockOsWriteFile
			err := WriteFile("someFile", []byte("someData"), os.ModePerm)
			if err != nil {
				if tc.wantErr == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				if tc.wantErr.Error() != err.Error() {
					t.Fatalf(`expected error "%v", got "%v"`, tc.wantErr, err)
				}
			} else {
				if tc.wantErr != nil {
					t.Fatalf(`expected error "%v", got nil`, tc.wantErr)
				}
			}
		})
	}
}

func TestCreateSymlink(t *testing.T) {
	testCases := []struct {
		name                  string
		mockOsCommandExecutor func() *mockOsCommandExecutor
		wantErr               error
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
			wantErr: errors.New("error when creating symlink from [src] to [dst]: some error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			osCommandExecutorProvider = tc.mockOsCommandExecutor()
			err := CreateSymlink("src", "dst")
			if err != nil {
				if tc.wantErr == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				if tc.wantErr.Error() != err.Error() {
					t.Fatalf(`expected error "%v", got "%v"`, tc.wantErr, err)
				}
			} else {
				if tc.wantErr != nil {
					t.Fatalf(`expected error "%v", got nil`, tc.wantErr)
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
	isDir bool
}

func (m *mockFileInfo) Name() string       { return "" }
func (m *mockFileInfo) Size() int64        { return 0 }
func (m *mockFileInfo) Mode() os.FileMode  { return 0 }
func (m *mockFileInfo) ModTime() time.Time { return time.Time{} }
func (m *mockFileInfo) IsDir() bool        { return m.isDir }
func (m *mockFileInfo) Sys() interface{}   { return nil }
