// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package dmg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/pkg/errors"
	"github.com/tiagomelo/macos-dmg-creator/validate"
)

const (
	contentsDir  = "Contents"
	macOsDir     = "Contents/MacOS"
	resourcesDir = "Contents/Resources"
	iconSetDir   = "icon.iconset"
)

// CreateParams is the input parameters for the Create function.
type CreateParams struct {
	// AppName is the name of the application.
	AppName string `validate:"required"`

	// AppBinaryPath is the path to the application binary.
	AppBinaryPath string `validate:"required"`

	// BundleIdentifier is the bundle identifier of the application.
	BundleIdentifier string `validate:"required"`

	// IconPath is the path to the icon file. Usable icon files are of type .png, .jpg, .gif, or .tiff.
	IconPath string `validate:"required"`

	// OutputDir is the directory where the DMG file will be created.
	OutputDir string `validate:"required"`
}

// Create creates a DMG file with the specified parameters.
func Create(params *CreateParams) (string, error) {
	// validate the input parameters.
	if err := validate.Check(params); err != nil {
		return "", errors.Wrap(err, "error when validating input parameters")
	}

	// temporary working directory for the application bundle.
	tmpWorkDir := filepath.Join(params.OutputDir, "tmp")
	if err := fsOpsProvider.MkdirAll(tmpWorkDir, os.ModePerm); err != nil {
		return "", errors.Wrap(err, "error when creating temp working directory")
	}
	// ensure the temporary working directory is cleaned up after use.
	defer func() {
		fsOpsProvider.DeleteDir(tmpWorkDir)
	}()

	appBundleSpinner := spinner.New(spinner.CharSets[14], 300*time.Millisecond)
	appBundleSpinner.Suffix = " creating application bundle..."
	appBundleSpinner.FinalMSG = "✔ creating application bundle...\n"
	appBundleSpinner.Start()

	// create the application bundle directory structure and files.
	createdAppBundleDirPath, err := createAppBundle(
		params.AppName,
		params.AppBinaryPath,
		params.IconPath,
		params.BundleIdentifier,
		tmpWorkDir,
	)
	appBundleSpinner.Stop()
	if err != nil {
		return "", errors.Wrap(err, "error when creating app bundle")
	}

	// create the DMG file from the application bundle.
	createdAppDmgPath, err := createAppDmg(createdAppBundleDirPath, tmpWorkDir, params.OutputDir)
	if err != nil {
		return "", errors.Wrap(err, "error when creating app DMG")
	}

	return createdAppDmgPath, nil
}

// createAppBundle creates the application bundle.
func createAppBundle(appName, appBinaryPath, iconPath, bundleIdentifier, outputDir string) (string, error) {
	appBundleDirName := fmt.Sprintf("%s.app", appName)
	appBundleDirPath := filepath.Join(outputDir, appBundleDirName)

	// create the application bundle directories.
	if err := createAppBundleDirectories(appBundleDirName, outputDir); err != nil {
		return "", errors.Wrap(err, "error when creating app bundle directories")
	}

	// create the icon set and copy the icons to the Resources directory.
	if err := createIconSet(iconPath, appBundleDirName, outputDir); err != nil {
		return "", errors.Wrap(err, "error when creating icon set")
	}

	// copy the application binary to the MacOS directory.
	if err := copyAppBinary(appBinaryPath, appBundleDirPath); err != nil {
		return "", errors.Wrap(err, "error when copying app binary")
	}

	// create the Info.plist file in the Resources directory.
	if err := createInfoPlistFile(appBinaryPath, appBundleDirName, outputDir, bundleIdentifier); err != nil {
		return "", errors.Wrap(err, "error when creating Info.plist file")
	}

	return appBundleDirPath, nil
}

// createAppBundleDirectories creates the necessary directories for the application bundle.
func createAppBundleDirectories(appBundleDirName, outputDir string) error {
	for _, dirName := range []string{
		filepath.Join(outputDir, appBundleDirName),
		filepath.Join(outputDir, iconSetDir),
		filepath.Join(outputDir, appBundleDirName, macOsDir),
		filepath.Join(outputDir, appBundleDirName, resourcesDir),
	} {
		err := fsOpsProvider.MkdirAll(dirName, os.ModePerm)
		if err != nil {
			return errors.Wrapf(err, "error when creating directory [%s]", dirName)
		}
	}
	return nil
}

// createIconSet creates the icon set directory structure.
func createIconSet(iconPath, appleBundleDirName, appBundleDirPath string) error {
	iconSizes := []int{16, 32, 64, 128, 256, 512, 1024}
	iconSetDirPath := filepath.Join(appBundleDirPath, iconSetDir)
	if err := sipsUtilityProvider.GenerateIcons(iconPath, iconSetDirPath, iconSizes...); err != nil {
		return errors.Wrap(err, "error when generating icons")
	}
	resourcesDirPath := filepath.Join(appBundleDirPath, appleBundleDirName, resourcesDir)
	if err := iconUtilProvider.GenerateIconSet(iconSetDirPath, resourcesDirPath); err != nil {
		return errors.Wrap(err, "error when generating icon set")
	}
	return nil
}

// copyAppBinary copies the application binary to the MacOS directory within the app bundle.
func copyAppBinary(appBinaryPath, appBundleDirPath string) error {
	macOsDirPath := filepath.Join(appBundleDirPath, macOsDir)
	if err := fsOpsProvider.CopyFile(appBinaryPath, macOsDirPath); err != nil {
		return errors.Wrapf(err, "error when copying file [%s] to [%s]", appBinaryPath, macOsDirPath)
	}
	return nil
}

// createInfoPlistFile creates the Info.plist file.
func createInfoPlistFile(appBinaryPath, appleBundleDirName, appBundleDirPath, bundleIdentifier string) error {
	infoPlist := strings.Replace(infoPlistTpl, "{{.AppName}}", filepath.Base(appBinaryPath), -1)
	infoPlist = strings.Replace(infoPlist, "{{.BundleIdentifier}}", bundleIdentifier, -1)
	contentsDirPath := filepath.Join(appBundleDirPath, appleBundleDirName, contentsDir)
	infoPlistPath := filepath.Join(contentsDirPath, "Info.plist")
	err := fsOpsProvider.WriteFile(infoPlistPath, []byte(infoPlist), os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "error when writing Info.plist file to [%s]", infoPlistPath)
	}
	return nil
}

// createAppDmg creates the DMG file for the application bundle.
func createAppDmg(appBundlePath, tmpWorkDir, outputDir string) (string, error) {
	dmgName := strings.TrimSuffix(filepath.Base(appBundlePath), ".app")
	if err := checkIfFinalDMGAlreadyExists(dmgName, outputDir); err != nil {
		return "", err
	}

	dmgTemplateSpinner := spinner.New(spinner.CharSets[14], 300*time.Millisecond)
	dmgTemplateSpinner.Suffix = " creating DMG template..."
	dmgTemplateSpinner.FinalMSG = "✔ creating DMG template...\n"
	dmgTemplateSpinner.Start()

	dmgTemplatePath, err := createDMGTemplate(dmgName, tmpWorkDir)
	dmgTemplateSpinner.Stop()
	if err != nil {
		return "", errors.Wrap(err, "error when creating DMG template")
	}

	mountDMGTemplateSpinner := spinner.New(spinner.CharSets[14], 300*time.Millisecond)
	mountDMGTemplateSpinner.Suffix = " mounting DMG template..."
	mountDMGTemplateSpinner.FinalMSG = "✔ mounting DMG template...\n"
	mountDMGTemplateSpinner.Start()

	mountPoint, err := mountDMGTemplate(dmgName, dmgTemplatePath)
	mountDMGTemplateSpinner.Stop()
	if err != nil {
		return "", errors.Wrap(err, "error when mounting DMG template")
	}

	setupDMGTemplateSpinner := spinner.New(spinner.CharSets[14], 300*time.Millisecond)
	setupDMGTemplateSpinner.Suffix = " setting up DMG template..."
	setupDMGTemplateSpinner.FinalMSG = "✔ setting up DMG template...\n"
	setupDMGTemplateSpinner.Start()

	err = setupDMGTemplate(mountPoint, appBundlePath)
	setupDMGTemplateSpinner.Stop()
	if err != nil {
		return "", errors.Wrap(err, "error when setting up DMG template")
	}

	unmountDMGTemplateSpinner := spinner.New(spinner.CharSets[14], 300*time.Millisecond)
	unmountDMGTemplateSpinner.Suffix = " unmounting DMG template..."
	unmountDMGTemplateSpinner.FinalMSG = "✔ unmounting DMG template...\n"
	unmountDMGTemplateSpinner.Start()

	err = unmountDMGTemplate(mountPoint)
	unmountDMGTemplateSpinner.Stop()
	if err != nil {
		return "", errors.Wrap(err, "error when unmounting DMG template")
	}

	return convertDmg(appBundlePath, dmgTemplatePath, outputDir)
}

// createDMGTemplate creates a DMG template for the application bundle.
func createDMGTemplate(dmgTemplateVolName, outputDir string) (string, error) {
	dmgTemplateFileName := fmt.Sprintf("%s-template.dmg", dmgTemplateVolName)
	dmgTemplatePath := filepath.Join(outputDir, dmgTemplateFileName)
	if err := hdiutilProvider.CreateDMG("100m", "APFS", dmgTemplateVolName, "GPTSPUD", dmgTemplatePath); err != nil {
		return "", err
	}
	return dmgTemplatePath, nil
}

// mountDMGTemplate mounts the DMG template to a temporary location.
func mountDMGTemplate(dmgTemplateVolName, dmgTemplatePath string) (string, error) {
	if err := hdiutilProvider.MoundDMG(dmgTemplateVolName, dmgTemplatePath); err != nil {
		return "", err
	}
	mountedDmgTemplatePath := fmt.Sprintf("/Volumes/%s", dmgTemplateVolName)
	return mountedDmgTemplatePath, nil
}

// setupDMGTemplate sets up the mounted DMG template with the application bundle.
func setupDMGTemplate(mountedDmgTemplatePath, createdAppBundleDirPath string) error {
	if err := createMacOsApplicationFolderSymlink(mountedDmgTemplatePath); err != nil {
		return errors.Wrap(err, "error when creating symlink for Applications folder")
	}
	if err := copyAppBundle(createdAppBundleDirPath, mountedDmgTemplatePath); err != nil {
		return errors.Wrapf(err, "error when copying app bundle to mounted DMG template at [%s]", mountedDmgTemplatePath)
	}
	return nil
}

// createMacOsApplicationFolderSymlink creates a symlink
// to the Applications folder in the mounted DMG template.
func createMacOsApplicationFolderSymlink(mountedDmgTemplatePath string) error {
	const symlinkName = "/Applications"
	if err := fsOpsProvider.CreateSymlink(symlinkName, mountedDmgTemplatePath); err != nil {
		return err
	}
	return nil
}

// copyAppBundle copies the application bundle to the mounted DMG template.
func copyAppBundle(createdAppBundleDirPath, mountedDmgTemplatePath string) error {
	if err := fsOpsProvider.CopyDir(createdAppBundleDirPath, mountedDmgTemplatePath); err != nil {
		return err
	}
	return nil
}

// unmountDMGTemplate unmounts the mounted DMG template.
func unmountDMGTemplate(mountedDmgTemplatePath string) error {
	if err := hdiutilProvider.UnmountDMG(mountedDmgTemplatePath); err != nil {
		return err
	}
	return nil
}

// checkIfFinalDMGAlreadyExists checks if the final DMG file
// already exists in the output directory.
func checkIfFinalDMGAlreadyExists(dmgName, outputDir string) error {
	dmgPath := filepath.Join(outputDir, fmt.Sprintf("%s.dmg", dmgName))
	exists, err := fsOpsProvider.FileExists(dmgPath)
	if err != nil {
		return errors.Wrapf(err, "error when checking if DMG file already exists at [%s]", dmgPath)
	}
	if exists {
		return errors.Errorf("DMG file already exists: [%s]", dmgPath)
	}
	return nil
}

// convertDmg converts the DMG template to the final DMG file.
func convertDmg(createdAppBundleDirPath, createdDmgTemplatePath, outputDir string) (string, error) {
	convertDMGSpinner := spinner.New(spinner.CharSets[14], 300*time.Millisecond)
	convertDMGSpinner.Suffix = " converting DMG template to final DMG..."
	convertDMGSpinner.FinalMSG = "✔ converting DMG template to final DMG...\n"
	convertDMGSpinner.Start()

	appDMGPath := filepath.Join(outputDir, fmt.Sprintf("%s.dmg", filepath.Base(strings.TrimSuffix(createdAppBundleDirPath, ".app"))))

	err := hdiutilProvider.ConvertDMG(createdDmgTemplatePath, appDMGPath)
	convertDMGSpinner.Stop()
	if err != nil {
		return "", errors.Wrap(err, "error when converting DMG template to final DMG")
	}
	return appDMGPath, nil
}
