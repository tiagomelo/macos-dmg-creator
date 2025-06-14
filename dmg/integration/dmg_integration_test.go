// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package dmg

import (
	"fmt"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/tiagomelo/macos-dmg-creator/dmg"
)

var createdDMGPath string

func TestMain(m *testing.M) {
	exitVal := m.Run()

	if createdDMGPath != "" {
		if err := os.Remove(createdDMGPath); err != nil {
			err = errors.Wrapf(err, "failed to remove created DMG at %s", createdDMGPath)
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// if err := os.Remove("dmg/integration/sampleapp/GreeterApp.dmg"); err != nil {
	// 	err = errors.Wrapf(err, "failed to remove sampleapp/GreeterApp.dmg")
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	os.Exit(exitVal)
}

func TestCreateDMG(t *testing.T) {
	const (
		appName       = "GreeterApp"
		appBinaryPath = "sampleapp/SampleApp"
		bundleID      = "info.tiago.greeterapp"
		iconPath      = "sampleapp/icon.png"
		outputDir     = "sampleapp"
	)

	var err error
	createdDMGPath, err = dmg.Create(&dmg.CreateParams{
		AppName:          appName,
		AppBinaryPath:    appBinaryPath,
		BundleIdentifier: bundleID,
		IconPath:         iconPath,
		OutputDir:        outputDir,
	})

	require.NoError(t, err)
	require.NotEmpty(t, createdDMGPath)
	require.FileExists(t, createdDMGPath)
}
