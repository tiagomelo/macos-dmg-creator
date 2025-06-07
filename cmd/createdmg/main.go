// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/tiagomelo/macos-dmg-creator/dmg"
)

// options defines the command line options for the program.
type options struct {
	AppName       string `long:"appName" description:"Application name" required:"true"`
	AppBinaryPath string `long:"appBinaryPath" description:"Path to the application binary" required:"true"`
	BundleID      string `long:"bundleIdentifier" description:"Bundle identifier for the application" required:"true"`
	IconPath      string `long:"iconPath" description:"Path to the application icon" required:"true"`
	OutputDir     string `long:"outputDir" description:"Directory to save the output DMG file" required:"true"`
}

func run(opts *options) error {
	createdDMGPath, err := dmg.Create(&dmg.CreateParams{
		AppName:          opts.AppName,
		AppBinaryPath:    opts.AppBinaryPath,
		BundleIdentifier: opts.BundleID,
		IconPath:         opts.IconPath,
		OutputDir:        opts.OutputDir,
	})
	if err != nil {
		return err
	}
	fmt.Println("\nDMG created successfully at:", createdDMGPath)
	return nil
}

func main() {
	var opts options
	parser := flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				fmt.Println(err)
				os.Exit(0)
			}
			fmt.Println(err)
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}
	if err := run(&opts); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
