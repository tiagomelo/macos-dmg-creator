// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package dmg

import "github.com/tiagomelo/macos-dmg-creator/iconutil"

// iconUtilProvider is a variable that holds the function
// that generates an icon set from a directory of icons.
var iconUtilProvider iconUtilOps = defaultIconUtil{}

// iconUtilOps defines an interface for generating an icon set.
type iconUtilOps interface {
	// GenerateIconSet generates an icon set from the specified icons directory.
	GenerateIconSet(iconsDir, outputDir string) error
}

// defaultIconUtil is the default implementation of iconUtil.
type defaultIconUtil struct{}

func (d defaultIconUtil) GenerateIconSet(iconsDir, outputDir string) error {
	return iconutil.GenerateIconSet(iconsDir, outputDir)
}
