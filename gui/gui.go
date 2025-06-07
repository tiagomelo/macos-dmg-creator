// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/tiagomelo/macos-dmg-creator/dmg"
)

// gui represents the graphical user interface for the app.
type gui struct {
	fyneApp    fyne.App
	fyneWindow fyne.Window
}

// New creates a new instance of the GUI for the app.
func New() *gui {
	const appTitle = "macOS DMG Creator"
	a := app.NewWithID("info.tiagomelo.macos-dmg-creator")
	w := a.NewWindow(appTitle)
	g := &gui{
		fyneApp:    a,
		fyneWindow: w,
	}
	g.setupUI()
	return g
}

// setupUI initializes the user interface components of the app.
func (g *gui) setupUI() {
	const (
		width  = 680
		height = 300
	)

	// set window dimensions.
	g.fyneWindow.Resize(fyne.NewSize(width, height))

	// setup the ui components.
	g.fyneWindow.SetContent(g.setupContent())

	// quit application when the window is closed.
	g.fyneWindow.SetCloseIntercept(func() {
		g.fyneApp.Quit()
	})
}

// setupContent creates and returns the main content of the GUI.
func (g *gui) setupContent() fyne.CanvasObject {
	const (
		dmgNameLabel               = "DMG name *"
		dmgNamePlaceholder         = "ExampleDMGName"
		dmgIconPathLabel           = "DMG icon path *"
		dmgIconPlaceholder         = "/path/to/dir/icon"
		dmgOutputLabel             = "DMG output path *"
		dmgOutputPlaceholder       = "/path/to/dir"
		dmgTemplatePathLabel       = "DMG template path"
		dmgTemplatePathPlaceholder = "/path/to/dir/template"
		appBinaryPathLabel         = "application binary path *"
		appBinaryPlaceholder       = "/path/to/dir"
		appBundleIDLabel           = "application bundle id *"
		appBundleIDPlaceholder     = "com.example.app"
		chooseLabel                = "choose..."
		requiredFielsLabel         = "* required fields"
	)

	// ==========================
	// DMG name
	// ==========================

	dmgNameEntry := widget.NewEntry()
	dmgNameEntry.SetPlaceHolder(dmgNamePlaceholder)
	dmgNameEntry.Validator = noSpaces

	// ==========================
	// Application binary path
	// ==========================

	appBinaryEntry := widget.NewEntry()
	appBinaryEntry.SetPlaceHolder(appBinaryPlaceholder)
	appBinaryEntry.Validator = noSpaces

	chooseAppBinaryPathButton := widget.NewButton(chooseLabel, func() {
		dialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, g.fyneWindow)
				return
			}
			if read == nil {
				return
			}
			appBinaryEntry.Text = read.URI().Path()
			appBinaryEntry.Refresh()
			appBinaryEntry.Validate()
		}, g.fyneWindow)

		dialog.Show()
	})
	chooseAppBinaryPathButton.Importance = widget.HighImportance

	// ==========================
	// DMG icon path
	// ==========================

	dmgIconEntry := widget.NewEntry()
	dmgIconEntry.SetPlaceHolder(dmgIconPlaceholder)
	dmgIconEntry.Validator = noSpaces

	chooseIconPathButton := widget.NewButton(chooseLabel, func() {
		dialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, g.fyneWindow)
				return
			}
			if read == nil {
				return
			}
			dmgIconEntry.Text = read.URI().Path()
			dmgIconEntry.Refresh()
			dmgIconEntry.Validate()
		}, g.fyneWindow)

		dialog.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg", ".gif", ".tiff"}))

		dialog.Show()
	})
	chooseIconPathButton.Importance = widget.HighImportance

	// ==========================
	// DMG output
	// ==========================

	dmgOutputEntry := widget.NewEntry()
	dmgOutputEntry.SetPlaceHolder(dmgOutputPlaceholder)
	dmgOutputEntry.Validator = noSpaces

	chooseDMGOutputPathButton := widget.NewButton(chooseLabel, func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, g.fyneWindow)
				return
			}
			if list == nil {
				return
			}
			dmgOutputEntry.Text = list.Path()
			dmgOutputEntry.Refresh()
			dmgOutputEntry.Validate()
		}, g.fyneWindow)
	})
	chooseDMGOutputPathButton.Importance = widget.HighImportance

	// ==========================
	// Application bundle id
	// ==========================

	appBundleIDEntry := widget.NewEntry()
	appBundleIDEntry.SetPlaceHolder(appBundleIDPlaceholder)
	appBundleIDEntry.Validator = noSpaces

	// ==========================
	// Progress bar dialog
	// ==========================

	progressBar := widget.NewProgressBarInfinite()
	rect := canvas.NewRectangle(color.Transparent)
	rect.SetMinSize(fyne.NewSize(200, 0))
	progressBarContainer := container.NewStack(rect, progressBar)
	progressBarDialog := dialog.NewCustom(
		"creating DMG...",
		"cancel",
		progressBarContainer,
		g.fyneWindow,
	)

	// ==========================
	// Form layout
	// ==========================

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: dmgNameLabel, Widget: dmgNameEntry},
			{Text: dmgIconPathLabel, Widget: dmgIconEntry},
			{Widget: chooseIconPathButton},
			{Text: appBinaryPathLabel, Widget: appBinaryEntry},
			{Widget: chooseAppBinaryPathButton},
			{Text: dmgOutputLabel, Widget: dmgOutputEntry},
			{Widget: chooseDMGOutputPathButton},
			{Text: appBundleIDLabel, Widget: appBundleIDEntry},
			{Widget: widget.NewLabelWithStyle(requiredFielsLabel, fyne.TextAlignCenter, fyne.TextStyle{Italic: true})},
		},
	}
	form.CancelText = "quit"
	form.OnCancel = func() {
		g.fyneApp.Quit()
	}
	form.SubmitText = "create DMG"
	form.OnSubmit = func() {
		form.Disable()
		progressBarDialog.Show()

		go func() {
			_, err := dmg.Create(&dmg.CreateParams{
				AppName:          dmgNameEntry.Text,
				AppBinaryPath:    appBinaryEntry.Text,
				BundleIdentifier: appBundleIDEntry.Text,
				IconPath:         dmgIconEntry.Text,
				OutputDir:        dmgOutputEntry.Text,
			})
			if err != nil {
				progressBarDialog.Hide()
				form.Enable()
				dialog.ShowError(err, g.fyneWindow)
				return
			}
			fyne.DoAndWait(func() {
				progressBarDialog.Hide()
				dialog.ShowInformation("Success", "DMG was successfully created!", g.fyneWindow)
				form.Enable()
			})
		}()
	}

	return container.NewVBox(form)
}

// Run starts the GUI application and enters the main event loop.
func (g *gui) Run() {
	g.fyneWindow.ShowAndRun()
}
