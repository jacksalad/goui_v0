package main

import (
	"goui/component"
	"goui/layout"
	"goui/render"
	"goui/window"
)

func main() {
	win, err := window.NewWindow(window.WindowConfig{
		Title:     "Simple Notepad",
		Width:     600,
		Height:    450,
		Resizable: false, // Keep it simple for now
	})
	if err != nil {
		panic(err)
	}

	// Setup Fonts
	uiFont := render.NewFont("Segoe UI", 14)
	editorFont := render.NewFont("Consolas", 16) // Monospace for editor

	// Root Layout (VBox)
	// We use a VBox to stack toolbar and editor
	vbox := &layout.VBoxLayout{
		Padding: 0,
		Spacing: 0,
	}
	win.Root.SetLayout(vbox)
	win.Root.BgColor = 0xFFCCCCCC

	// Toolbar
	toolbar := component.NewPanel(0, 0, 600, 40)
	toolbar.BgColor = 0xFFDDDDDD
	toolbar.SetLayout(&layout.HBoxLayout{
		Padding: 5,
		Spacing: 10,
	})
	win.Root.Add(toolbar)

	// Toolbar Buttons
	btnNew := component.NewButton("Clear")
	btnNew.Font = uiFont
	toolbar.Add(btnNew)

	btnInfo := component.NewButton("About")
	btnInfo.Font = uiFont
	toolbar.Add(btnInfo)

	// Editor Area
	// We want it to fill the rest of the window.
	// Since VBoxLayout stacks, we just need to size it correctly.
	// Window Height (450) - Toolbar Height (40) = 410.
	editor := component.NewTextArea(600, 410)
	editor.Font = editorFont
	editor.Text = "Welcome to Simple Notepad!\nStart typing here..."

	win.Root.Add(editor)

	// Wire up events
	btnNew.OnClick = func() {
		editor.Text = ""
		editor.RequestRepaint()
		win.SetFocus(editor)
	}

	btnInfo.OnClick = func() {
		editor.Text += "\n\n[About]\nSimple Notepad built with GoUI.\nSupports basic editing."
		editor.RequestRepaint()
	}

	win.Show()
	win.Run()
}
