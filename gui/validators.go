// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package gui

import (
	"errors"
	"strings"
)

// notEmpty is a simple validator that checks if the input string is empty.
func notEmpty(s string) error {
	if s == "" {
		return errors.New("cannot be blank")
	}
	return nil
}

// noSpaces checks that the input does not contain any space characters.
func noSpaces(input string) error {
	if strings.Contains(input, " ") {
		return errors.New("cannot contain spaces")
	}
	return notEmpty(input)
}

// optionalNotEmpty checks if the input is empty, and if not, applies the notEmpty validator.
func optionalNoSpaces(input string) error {
	if input == "" {
		return nil // Allow empty input
	}
	return noSpaces(input)
}
