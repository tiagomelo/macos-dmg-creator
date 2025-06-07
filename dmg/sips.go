// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package dmg

import "github.com/tiagomelo/macos-dmg-creator/sips"

// sipsUtilityProvider is a variable that holds the function
// that generates icons using the sips utility.
var sipsUtilityProvider sipsOps = defaultSips{}

// sipsOps defines an interface for generating icons using the sips utility.
type sipsOps interface {
	// GenerateIcons generates icons from the specified icon file at the given sizes.
	GenerateIcons(iconPath, outputDir string, sizes ...int) error
}

// defaultSips is the default implementation of sipsUtility.
type defaultSips struct{}

func (d defaultSips) GenerateIcons(iconPath, outputDir string, sizes ...int) error {
	return sips.GenerateIcons(iconPath, outputDir, sizes...)
}
