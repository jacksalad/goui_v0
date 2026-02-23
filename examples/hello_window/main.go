package main

import (
	"fmt"
	"goui/event"
	"goui/window"
	"log"
)

func main() {
	w, err := window.NewWindow(window.WindowConfig{
		Title:     "Hello Goui",
		Width:     800,
		Height:    600,
		Resizable: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Subscribe to close event
	closeCh := w.EventBus.Subscribe(event.EventClose)

	go func() {
		for evt := range closeCh {
			if evt.Type == event.EventClose {
				fmt.Println("Window closing...")
			}
		}
	}()

	w.Show()
	w.Run()
}
