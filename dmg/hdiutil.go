// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package dmg

import "github.com/tiagomelo/macos-dmg-creator/hdiutil"

// hdiutilProvider is a variable that holds the function
// that interacts with the hdiutil command-line tool.
var hdiutilProvider hdiutilOps = defaultHdiutil{}

// hdiutilOps defines an interface for interacting with the hdiutilOps command-line tool.
type hdiutilOps interface {
	// CreateDMG creates a DMG file from the specified source directory.
	CreateDMG(size, fs, volName, layout, output string) error

	// MoundDMG mounts the DMG file and returns the mount point.
	MoundDMG(dmgVolName, dmgPath string) error

	// UnmountDMG unmounts the DMG file that was previously mounted.
	UnmountDMG(volName string) error

	// ConvertDMG converts a DMG file to a different format.
	ConvertDMG(dmgPath, dmgOutputFileName string) error
}

// defaultHdiutil is the default implementation of hdiutilOps.
type defaultHdiutil struct{}

func (d defaultHdiutil) CreateDMG(size, fs, volName, layout, output string) error {
	return hdiutil.CreateDMG(size, fs, volName, layout, output)
}

func (d defaultHdiutil) MoundDMG(dmgVolName, dmgPath string) error {
	return hdiutil.MountDMG(dmgVolName, dmgPath)
}

func (d defaultHdiutil) UnmountDMG(volName string) error {
	return hdiutil.UnmountDMG(volName)
}

func (d defaultHdiutil) ConvertDMG(dmgPath, dmgOutputFileName string) error {
	return hdiutil.ConvertDMG(dmgPath, dmgOutputFileName)
}
