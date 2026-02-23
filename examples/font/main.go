package main

import (
	"runtime"

	"github.com/jacksalad/goui_v0/component"
	"github.com/jacksalad/goui_v0/layout"
	"github.com/jacksalad/goui_v0/render"
	"github.com/jacksalad/goui_v0/window"
)

func main() {
	runtime.LockOSThread()

	// 1. Create a window
	win, err := window.NewWindow(window.WindowConfig{
		Title:     "Font Example - Microsoft YaHei & Consolas",
		Width:     800,
		Height:    600,
		Resizable: true,
	})
	if err != nil {
		panic(err)
	}

	// 2. Load a system font (or local font)
	// Example: Loading "Microsoft YaHei" (common on Windows)
	font := render.NewFont("Microsoft YaHei", 12)
	win.Renderer.SetFont(font)
	defer font.Close()

	// 3. Create components
	// Use win.Root directly
	win.Root.SetLayout(&layout.VBoxLayout{Spacing: 10, Padding: 20})

	label := component.NewLabel("This is Microsoft YaHei (12pt)")
	win.Root.Add(label)

	tb := component.NewTextBox(300)
	tb.Text = "Supports Chinese Input: 你好，世界！"
	win.Root.Add(tb)

	btn := component.NewButton("Change Font to Consolas (16pt)")
	btn.OnClick = func() {
		// Note: This creates a new font every click, should ideally cache it or close old one
		// But for example it's fine.
		// Also we should close the old font if we are replacing it, but here we just leak it for simplicity
		// or we can keep track of it.
		newFont := render.NewFont("Consolas", 16)
		win.Renderer.SetFont(newFont)

		label.Text = "Now using Consolas (16pt)"
		btn.Text = "Font Changed!"
	}
	win.Root.Add(btn)

	// 4. Run the message loop
	win.Show()
	win.Run()
}
