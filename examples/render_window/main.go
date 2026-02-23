package main

import (
	"github.com/jacksalad/goui_v0/window"
	"log"
)

func main() {
	w, err := window.NewWindow(window.WindowConfig{
		Title:     "Render Window",
		Width:     800,
		Height:    600,
		Resizable: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	
	// Simple loop to draw
	go func() {
		for {
			// Draw something
			canvas := w.Renderer.BeginFrame()
			canvas.Clear(0xFF0000FF) // Blue background (ARGB)
			canvas.FillRect(100, 100, 200, 150, 0xFFFF0000) // Red rect
			w.Renderer.EndFrame()
			w.Renderer.Present()
			
			// We can't easily sleep without importing time
			// and busy loop is bad, but for test it's fine.
		}
	}()
	
	w.Show()
	w.Run()
}
