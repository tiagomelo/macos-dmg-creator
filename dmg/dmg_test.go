// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package dmg

import (
	"os"
	"testing"

	"github.com/pkg/errors"
)

func TestCreate(t *testing.T) {
	testCases := []struct {
		name                    string
		params                  *CreateParams
		mockFsOpsProvider       func() *mockFsOpsProvider
		mockSipsUtilityProvider func() *mockSipsUtilityProvider
		mockIconUtilProvider    func() *mockIconUtilProvider
		mockHdiutilProvider     func() *mockHdiutilProvider
		want                    string
		wantErr                 error
	}{
		{
			name: "happy path",
			params: &CreateParams{
				AppName:          "testAppName",
				AppBinaryPath:    "testAppBinaryPath",
				BundleIdentifier: "testBundleIdentifier",
				IconPath:         "testIconPath",
				OutputDir:        "outputDir",
			},
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{}
			},
			mockSipsUtilityProvider: func() *mockSipsUtilityProvider {
				return &mockSipsUtilityProvider{}
			},
			mockIconUtilProvider: func() *mockIconUtilProvider {
				return &mockIconUtilProvider{}
			},
			mockHdiutilProvider: func() *mockHdiutilProvider {
				return &mockHdiutilProvider{}
			},
			want: "outputDir/testAppName.dmg",
		},
		{
			name: "error when validating input parameters",
			params: &CreateParams{
				AppName:          "",
				AppBinaryPath:    "testAppBinaryPath",
				BundleIdentifier: "testBundleIdentifier",
				IconPath:         "testIconPath",
				OutputDir:        "outputDir",
			},
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{}
			},
			mockSipsUtilityProvider: func() *mockSipsUtilityProvider {
				return &mockSipsUtilityProvider{}
			},
			mockIconUtilProvider: func() *mockIconUtilProvider {
				return &mockIconUtilProvider{}
			},
			mockHdiutilProvider: func() *mockHdiutilProvider {
				return &mockHdiutilProvider{}
			},
			wantErr: errors.New("error when validating input parameters: AppName: AppName is a required field"),
		},
		{
			name: "error when creating temp working directory",
			params: &CreateParams{
				AppName:          "testAppName",
				AppBinaryPath:    "testAppBinaryPath",
				BundleIdentifier: "testBundleIdentifier",
				IconPath:         "testIconPath",
				OutputDir:        "outputDir",
			},
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{
					expectedMkdirAllErr: os.ErrPermission,
				}
			},
			mockSipsUtilityProvider: func() *mockSipsUtilityProvider {
				return &mockSipsUtilityProvider{}
			},
			mockIconUtilProvider: func() *mockIconUtilProvider {
				return &mockIconUtilProvider{}
			},
			mockHdiutilProvider: func() *mockHdiutilProvider {
				return &mockHdiutilProvider{}
			},
			wantErr: errors.Wrap(os.ErrPermission, "error when creating temp working directory"),
		},
		{
			name: "error when creating app bundle",
			params: &CreateParams{
				AppName:          "testAppName",
				AppBinaryPath:    "testAppBinaryPath",
				BundleIdentifier: "testBundleIdentifier",
				IconPath:         "testIconPath",
				OutputDir:        "outputDir",
			},
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{}
			},
			mockSipsUtilityProvider: func() *mockSipsUtilityProvider {
				return &mockSipsUtilityProvider{
					expectedGenerateIconsErr: os.ErrPermission,
				}
			},
			mockIconUtilProvider: func() *mockIconUtilProvider {
				return &mockIconUtilProvider{}
			},
			mockHdiutilProvider: func() *mockHdiutilProvider {
				return &mockHdiutilProvider{}
			},
			wantErr: errors.Wrap(os.ErrPermission, "error when creating app bundle: error when creating icon set: error when generating icons"),
		},
		{
			name: "error when creating app DMG",
			params: &CreateParams{
				AppName:          "testAppName",
				AppBinaryPath:    "testAppBinaryPath",
				BundleIdentifier: "testBundleIdentifier",
				IconPath:         "testIconPath",
				OutputDir:        "outputDir",
			},
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{}
			},
			mockSipsUtilityProvider: func() *mockSipsUtilityProvider {
				return &mockSipsUtilityProvider{}
			},
			mockIconUtilProvider: func() *mockIconUtilProvider {
				return &mockIconUtilProvider{}
			},
			mockHdiutilProvider: func() *mockHdiutilProvider {
				return &mockHdiutilProvider{
					expectedMountDMGErr: os.ErrPermission,
				}
			},
			wantErr: errors.Wrap(os.ErrPermission, "error when creating app DMG: error when mounting DMG template"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fsOpsProvider = tc.mockFsOpsProvider()
			sipsUtilityProvider = tc.mockSipsUtilityProvider()
			iconUtilProvider = tc.mockIconUtilProvider()
			hdiutilProvider = tc.mockHdiutilProvider()

			got, err := Create(tc.params)
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
			if got != tc.want {
				t.Fatalf(`expected DMG file path "%s", got "%s"`, tc.want, got)
			}

		})
	}
}

func Test_createAppBundle(t *testing.T) {
	testCases := []struct {
		name                    string
		mockFsOpsProvider       func() *mockFsOpsProvider
		mockSipsUtilityProvider func() *mockSipsUtilityProvider
		mockIconUtilProvider    func() *mockIconUtilProvider
		want                    string
		wantErr                 error
	}{
		{
			name: "happy path",
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{}
			},
			mockSipsUtilityProvider: func() *mockSipsUtilityProvider {
				return &mockSipsUtilityProvider{}
			},
			mockIconUtilProvider: func() *mockIconUtilProvider {
				return &mockIconUtilProvider{}
			},
			want: "testOutputDir/testAppName.app",
		},
		{
			name: "error when creating directories",
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{
					expectedMkdirAllErr: os.ErrPermission,
				}
			},
			mockSipsUtilityProvider: func() *mockSipsUtilityProvider {
				return &mockSipsUtilityProvider{}
			},
			mockIconUtilProvider: func() *mockIconUtilProvider {
				return &mockIconUtilProvider{}
			},
			wantErr: errors.Wrap(os.ErrPermission, "error when creating app bundle directories: error when creating directory [testOutputDir/testAppName.app]"),
		},
		{
			name: "error when creating icon set",
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{}
			},
			mockSipsUtilityProvider: func() *mockSipsUtilityProvider {
				return &mockSipsUtilityProvider{
					expectedGenerateIconsErr: os.ErrPermission,
				}
			},
			mockIconUtilProvider: func() *mockIconUtilProvider {
				return &mockIconUtilProvider{}
			},
			wantErr: errors.Wrap(os.ErrPermission, "error when creating icon set: error when generating icons"),
		},
		{
			name: "error when copying app binary",
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{
					expectedCopyFileErr: os.ErrPermission,
				}
			},
			mockSipsUtilityProvider: func() *mockSipsUtilityProvider {
				return &mockSipsUtilityProvider{}
			},
			mockIconUtilProvider: func() *mockIconUtilProvider {
				return &mockIconUtilProvider{}
			},
			wantErr: errors.Wrap(os.ErrPermission, "error when copying app binary: error when copying file [testAppBinaryPath] to [testOutputDir/testAppName.app/Contents/MacOS]"),
		},
		{
			name: "error when creating Info.plist file",
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{
					expectedWriteFileErr: os.ErrPermission,
				}
			},
			mockSipsUtilityProvider: func() *mockSipsUtilityProvider {
				return &mockSipsUtilityProvider{}
			},
			mockIconUtilProvider: func() *mockIconUtilProvider {
				return &mockIconUtilProvider{}
			},
			wantErr: errors.Wrap(os.ErrPermission, "error when creating Info.plist file: error when writing Info.plist file to [testOutputDir/testAppName.app/Contents/Info.plist]"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fsOpsProvider = tc.mockFsOpsProvider()
			sipsUtilityProvider = tc.mockSipsUtilityProvider()
			iconUtilProvider = tc.mockIconUtilProvider()

			got, err := createAppBundle(
				"testAppName",
				"testAppBinaryPath",
				"testIconPath",
				"testBundleIdentifier",
				"testOutputDir",
			)

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

			if got != tc.want {
				t.Fatalf(`expected app bundle directory path "%s", got "%s"`, tc.want, got)
			}
		})
	}
}

func Test_createAppBundleDirectories(t *testing.T) {
	testCases := []struct {
		name              string
		mockFsOpsProvider func() *mockFsOpsProvider
		wantErr           error
	}{
		{
			name: "happy path",
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{}
			},
		},
		{
			name: "error when creating directories",
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{
					expectedMkdirAllErr: os.ErrPermission,
				}
			},
			wantErr: errors.Wrap(os.ErrPermission, "error when creating directory [testOutputDir/testAppBundleDirName]"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fsOpsProvider = tc.mockFsOpsProvider()

			err := createAppBundleDirectories("testAppBundleDirName", "testOutputDir")

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

func Test_createIconSet(t *testing.T) {
	testCases := []struct {
		name                    string
		mockSipsUtilityProvider func() *mockSipsUtilityProvider
		mockIconUtilProvider    func() *mockIconUtilProvider
		wantErr                 error
	}{
		{
			name: "happy path",
			mockSipsUtilityProvider: func() *mockSipsUtilityProvider {
				return &mockSipsUtilityProvider{}
			},
			mockIconUtilProvider: func() *mockIconUtilProvider {
				return &mockIconUtilProvider{}
			},
		},
		{
			name: "error when generating icons",
			mockSipsUtilityProvider: func() *mockSipsUtilityProvider {
				return &mockSipsUtilityProvider{
					expectedGenerateIconsErr: os.ErrPermission,
				}
			},
			mockIconUtilProvider: func() *mockIconUtilProvider {
				return &mockIconUtilProvider{}
			},
			wantErr: errors.Wrap(os.ErrPermission, "error when generating icons"),
		},
		{
			name: "error when generating icon set",
			mockSipsUtilityProvider: func() *mockSipsUtilityProvider {
				return &mockSipsUtilityProvider{}
			},
			mockIconUtilProvider: func() *mockIconUtilProvider {
				return &mockIconUtilProvider{
					expectedGenerateIconSetErr: os.ErrPermission,
				}
			},
			wantErr: errors.Wrap(os.ErrPermission, "error when generating icon set"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sipsUtilityProvider = tc.mockSipsUtilityProvider()
			iconUtilProvider = tc.mockIconUtilProvider()

			err := createIconSet(
				"testIconPath",
				"testAppBundleDirName",
				"testAppBundleDirPath",
			)

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

func Test_copyAppBinary(t *testing.T) {
	testCases := []struct {
		name              string
		mockFsOpsProvider func() *mockFsOpsProvider
		wantErr           error
	}{
		{
			name: "happy path",
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{}
			},
		},
		{
			name: "error when copying app binary",
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{
					expectedCopyFileErr: os.ErrPermission,
				}
			},
			wantErr: errors.Wrap(os.ErrPermission, "error when copying file [testAppBinaryPath] to [testOutputDir/Contents/MacOS]"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fsOpsProvider = tc.mockFsOpsProvider()

			err := copyAppBinary("testAppBinaryPath", "testOutputDir")
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

func Test_createInfoPlistFile(t *testing.T) {
	testCases := []struct {
		name              string
		mockFsOpsProvider func() *mockFsOpsProvider
		wantErr           error
	}{
		{
			name: "happy path",
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{}
			},
		},
		{
			name: "error when writing file",
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{
					expectedWriteFileErr: os.ErrPermission,
				}
			},
			wantErr: errors.Wrap(os.ErrPermission, "error when writing Info.plist file to [testIconPath/testBundleIdentifier/Contents/Info.plist]"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fsOpsProvider = tc.mockFsOpsProvider()

			err := createInfoPlistFile(
				"testAppName",
				"testBundleIdentifier",
				"testIconPath",
				"testOutputDir",
			)

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

func Test_createAppDmg(t *testing.T) {
	testCases := []struct {
		name                string
		mockFsOpsProvider   func() *mockFsOpsProvider
		mockHdiutilProvider func() *mockHdiutilProvider
		want                string
		wantErr             error
	}{
		{
			name: "happy path",
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{}
			},
			mockHdiutilProvider: func() *mockHdiutilProvider {
				return &mockHdiutilProvider{}
			},
			want: "outputDir/testAppBundleDirPath.dmg",
		},
		{
			name: "error checking if final DMG already exists",
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{
					expectedFileExistsErr: os.ErrPermission,
				}
			},
			mockHdiutilProvider: func() *mockHdiutilProvider {
				return &mockHdiutilProvider{}
			},
			wantErr: errors.Wrap(os.ErrPermission, "error when checking if DMG file already exists at [outputDir/testAppBundleDirPath.dmg]"),
		},
		{
			name: "error creating DMG template",
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{}
			},
			mockHdiutilProvider: func() *mockHdiutilProvider {
				return &mockHdiutilProvider{
					expectedCreateDMGErr: os.ErrPermission,
				}
			},
			wantErr: errors.Wrap(os.ErrPermission, "error when creating DMG template"),
		},
		{
			name: "error mounting DMG template",
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{}
			},
			mockHdiutilProvider: func() *mockHdiutilProvider {
				return &mockHdiutilProvider{
					expectedMountDMGErr: os.ErrPermission,
				}
			},
			wantErr: errors.Wrap(os.ErrPermission, "error when mounting DMG template"),
		},
		{
			name: "error setting up DMG template",
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{
					expectedCreateSymlinkErr: os.ErrPermission,
				}
			},
			mockHdiutilProvider: func() *mockHdiutilProvider {
				return &mockHdiutilProvider{}
			},
			wantErr: errors.New("error when setting up DMG template: error when creating symlink for Applications folder: permission denied"),
		},
		{
			name: "error unmounting DMG template",
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{}
			},
			mockHdiutilProvider: func() *mockHdiutilProvider {
				return &mockHdiutilProvider{
					expectedUnmountDMGErr: os.ErrPermission,
				}
			},
			wantErr: errors.Wrap(os.ErrPermission, "error when unmounting DMG template"),
		},
		{
			name: "error converting DMG",
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{}
			},
			mockHdiutilProvider: func() *mockHdiutilProvider {
				return &mockHdiutilProvider{
					expectedConvertDMGErr: os.ErrPermission,
				}
			},
			wantErr: errors.Wrap(os.ErrPermission, "error when converting DMG template to final DMG"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fsOpsProvider = tc.mockFsOpsProvider()
			hdiutilProvider = tc.mockHdiutilProvider()

			got, err := createAppDmg(
				"testAppBundleDirPath",
				"tmpWorkDir",
				"outputDir",
			)

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
			if got != tc.want {
				t.Fatalf(`expected DMG file path "%s", got "%s"`, tc.want, got)
			}
		})
	}
}

func Test_checkIfFinalDMGAlreadyExists(t *testing.T) {
	testCases := []struct {
		name              string
		mockFsOpsProvider func() *mockFsOpsProvider
		wantErr           error
	}{
		{
			name: "does not exists",
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{}
			},
		},
		{
			name: "exists",
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{
					expectedFileExists: true,
				}
			},
			wantErr: errors.New("DMG file already exists: [outputDir/dmgName.dmg]"),
		},
		{
			name: "error checking file existence",
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{
					expectedFileExistsErr: os.ErrPermission,
				}
			},
			wantErr: errors.Wrap(os.ErrPermission, "error when checking if DMG file already exists at [outputDir/dmgName.dmg]"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fsOpsProvider = tc.mockFsOpsProvider()

			err := checkIfFinalDMGAlreadyExists(
				"dmgName",
				"outputDir",
			)

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

func Test_createDMGTemplate(t *testing.T) {
	testCases := []struct {
		name                string
		mockHdiutilProvider func() *mockHdiutilProvider
		want                string
		wantErr             error
	}{
		{
			name: "happy path",
			mockHdiutilProvider: func() *mockHdiutilProvider {
				return &mockHdiutilProvider{}
			},
			want: "outputDir/dmgTemplateVolName-template.dmg",
		},
		{
			name: "error creating DMG",
			mockHdiutilProvider: func() *mockHdiutilProvider {
				return &mockHdiutilProvider{
					expectedCreateDMGErr: os.ErrPermission,
				}
			},
			wantErr: os.ErrPermission,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hdiutilProvider = tc.mockHdiutilProvider()

			got, err := createDMGTemplate(
				"dmgTemplateVolName",
				"outputDir",
			)

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

			if got != tc.want {
				t.Fatalf(`expected app bundle directory path "%s", got "%s"`, tc.want, got)
			}
		})
	}
}

func Test_mountDMGTemplate(t *testing.T) {
	testCases := []struct {
		name                string
		mockHdiutilProvider func() *mockHdiutilProvider
		want                string
		wantErr             error
	}{
		{
			name: "happy path",
			mockHdiutilProvider: func() *mockHdiutilProvider {
				return &mockHdiutilProvider{}
			},
			want: "/Volumes/dmgTemplateVolName",
		},
		{
			name: "error mounting DMG",
			mockHdiutilProvider: func() *mockHdiutilProvider {
				return &mockHdiutilProvider{
					expectedMountDMGErr: os.ErrPermission,
				}
			},
			wantErr: os.ErrPermission,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hdiutilProvider = tc.mockHdiutilProvider()

			got, err := mountDMGTemplate(
				"dmgTemplateVolName",
				"outputDir/dmgTemplateVolName-template.dmg",
			)

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

			if got != tc.want {
				t.Fatalf(`expected mounted DMG path "%s", got "%s"`, tc.want, got)
			}
		})
	}
}

func Test_setupDMGTemplate(t *testing.T) {
	testCases := []struct {
		name              string
		mockFsOpsProvider func() *mockFsOpsProvider
		wantErr           error
	}{
		{
			name: "happy path",
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{}
			},
		},
		{
			name: "error creating symlink for Applications folder",
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{
					expectedCreateSymlinkErr: os.ErrPermission,
				}
			},
			wantErr: errors.Wrap(os.ErrPermission, "error when creating symlink for Applications folder"),
		},
		{
			name: "error copying app bundle to mounted DMG template",
			mockFsOpsProvider: func() *mockFsOpsProvider {
				return &mockFsOpsProvider{
					expectedCopyDirErr: os.ErrPermission,
				}
			},
			wantErr: errors.Wrap(os.ErrPermission, "error when copying app bundle to mounted DMG template at [/Volumes/dmgTemplateVolName]"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fsOpsProvider = tc.mockFsOpsProvider()

			err := setupDMGTemplate(
				"/Volumes/dmgTemplateVolName",
				"testAppBundleDirPath",
			)

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

func Test_unmountDMGTemplate(t *testing.T) {
	testCases := []struct {
		name                string
		mockHdiutilProvider func() *mockHdiutilProvider
		wantErr             error
	}{
		{
			name: "happy path",
			mockHdiutilProvider: func() *mockHdiutilProvider {
				return &mockHdiutilProvider{}
			},
		},
		{
			name: "error unmounting DMG",
			mockHdiutilProvider: func() *mockHdiutilProvider {
				return &mockHdiutilProvider{
					expectedUnmountDMGErr: os.ErrPermission,
				}
			},
			wantErr: os.ErrPermission,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hdiutilProvider = tc.mockHdiutilProvider()

			err := unmountDMGTemplate("/Volumes/dmgTemplateVolName")
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

func Test_convertDmg(t *testing.T) {
	testCases := []struct {
		name                string
		mockHdiutilProvider func() *mockHdiutilProvider
		want                string
		wantErr             error
	}{
		{
			name: "happy path",
			mockHdiutilProvider: func() *mockHdiutilProvider {
				return &mockHdiutilProvider{}
			},
			want: "outputDir/testAppBundleDirPath.dmg",
		},
		{
			name: "error converting DMG",
			mockHdiutilProvider: func() *mockHdiutilProvider {
				return &mockHdiutilProvider{
					expectedConvertDMGErr: os.ErrPermission,
				}
			},
			wantErr: errors.Wrap(os.ErrPermission, "error when converting DMG template to final DMG"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hdiutilProvider = tc.mockHdiutilProvider()

			got, err := convertDmg(
				"outputDir/testAppBundleDirPath.app",
				"outputDir/tmp/dmgTemplateVolName-template.dmg",
				"outputDir",
			)

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

			if got != tc.want {
				t.Fatalf(`expected converted DMG path "%s", got "%s"`, tc.want, got)
			}
		})
	}
}

type mockFsOpsProvider struct {
	expectedDirExists              bool
	expectedDirExistsErr           error
	expectedVolumeDoesNotExists    bool
	expectedVolumeDoesNotExistsErr error
	expectedFileExists             bool
	expectedFileExistsErr          error
	expectedMkdirAllErr            error
	expectedCopyFileErr            error
	expectedCopyDirErr             error
	expectedWriteFileErr           error
	expectedCreateSymlinkErr       error
	expectedDeleteDirErr           error
}

func (m *mockFsOpsProvider) DirExists(path string) (bool, error) {
	return m.expectedDirExists, m.expectedDirExistsErr
}

func (m *mockFsOpsProvider) VolumeDoesNotExist(path string) (bool, error) {
	return m.expectedVolumeDoesNotExists, m.expectedVolumeDoesNotExistsErr
}

func (m *mockFsOpsProvider) MkdirAll(path string, perm os.FileMode) error {
	return m.expectedMkdirAllErr
}

func (m *mockFsOpsProvider) CopyFile(src, dst string) error {
	return m.expectedCopyFileErr
}

func (m *mockFsOpsProvider) CopyDir(src, dst string) error {
	return m.expectedCopyDirErr
}

func (m *mockFsOpsProvider) DeleteDir(path string) error {
	return m.expectedDeleteDirErr
}

func (m *mockFsOpsProvider) WriteFile(name string, data []byte, perm os.FileMode) error {
	return m.expectedWriteFileErr
}

func (m *mockFsOpsProvider) CreateSymlink(src, dst string) error {
	return m.expectedCreateSymlinkErr
}

func (m *mockFsOpsProvider) FileExists(path string) (bool, error) {
	return m.expectedFileExists, m.expectedFileExistsErr
}

type mockSipsUtilityProvider struct {
	expectedGenerateIconsErr error
}

func (m *mockSipsUtilityProvider) GenerateIcons(iconPath, iconSetDirPath string, sizes ...int) error {
	return m.expectedGenerateIconsErr
}

type mockIconUtilProvider struct {
	expectedGenerateIconSetErr error
}

func (m *mockIconUtilProvider) GenerateIconSet(iconSetDirPath, resourcesDirPath string) error {
	return m.expectedGenerateIconSetErr
}

type mockHdiutilProvider struct {
	expectedCreateDMGErr  error
	expectedConvertDMGErr error
	expectedMountDMGErr   error
	expectedUnmountDMGErr error
}

func (m *mockHdiutilProvider) CreateDMG(size, fs, volName, layout, output string) error {
	return m.expectedCreateDMGErr
}

func (m *mockHdiutilProvider) ConvertDMG(dmgPath, dmgOutputFileName string) error {
	return m.expectedConvertDMGErr
}

func (m *mockHdiutilProvider) MoundDMG(dmgVolName, dmgPath string) error {
	return m.expectedMountDMGErr
}

func (m *mockHdiutilProvider) UnmountDMG(volName string) error {
	return m.expectedUnmountDMGErr
}
