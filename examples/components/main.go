package main

import (
	"fmt"
	"log"

	"github.com/jacksalad/goui_v0/component"
	"github.com/jacksalad/goui_v0/window"
)

func main() {
	// Create main window
	w, err := window.NewWindow(window.WindowConfig{
		Title:     "Components Demo",
		Width:     800,
		Height:    600,
		Resizable: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a main panel
	mainPanel := component.NewPanel(0, 0, 800, 600)
	mainPanel.BgColor = 0xFFFFFFFF // White
	w.Add(mainPanel)

	// Create a sub-panel
	subPanel := component.NewPanel(50, 50, 700, 500)
	subPanel.BgColor = 0xFFF0F0F0 // Light gray
	mainPanel.Add(subPanel)

	// Create a label
	label := component.NewLabel("Hello from GDI!")
	label.SetBounds(100, 100, 200, 30)
	subPanel.Add(label)

	// Create a button
	btn := component.NewButton("Click Me")
	btn.SetBounds(100, 150, 120, 40)
	btn.OnClick = func() {
		fmt.Println("Button Clicked!")
		label.Text = "Button Clicked!"
		// The window will repaint automatically on the next event or immediately
	}
	subPanel.Add(btn)

	w.Show()
	w.Run()
}
