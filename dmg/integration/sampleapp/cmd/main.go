// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	const (
		appTitle = "greeter"
		width    = 300
		height   = 180
	)

	greeterApp := app.NewWithID("info.tiagomelo.greeter")
	window := greeterApp.NewWindow(appTitle)

	window.Resize(fyne.NewSize(width, height))

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Enter your name")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: nameEntry},
		},
		OnSubmit: func() {
			name := nameEntry.Text
			if name == "" {
				name = "World"
			}
			dialog.ShowInformation("Greeting", "Hello, "+name+"!", window)
		},
		OnCancel: func() {
			greeterApp.Quit()
		},
	}

	window.SetContent(container.NewVBox(form))
	window.SetCloseIntercept(func() {
		greeterApp.Quit()
	})

	window.ShowAndRun()
}
