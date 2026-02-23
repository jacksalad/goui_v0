# Getting Started with GoUI

## Prerequisites

*   **Go 1.18+**: Ensure you have Go installed.
*   **Windows**: GoUI currently supports Windows (uses `user32.dll` and `gdi32.dll`).

## Installation

Since GoUI is a local library in this context, you don't need to `go get` it if you are working within the repository.

If you were to use it as a module:

```bash
go get github.com/jacksalad/goui_v0
```

## Hello World

Here is a minimal example to create a window with a button.

```go
package main

import (
	"fmt"
	"github.com/jacksalad/goui_v0/component"
	"github.com/jacksalad/goui_v0/layout"
	"github.com/jacksalad/goui_v0/window"
)

func main() {
	// 1. Create a new Window
	win, err := window.NewWindow(window.WindowConfig{
		Title:     "Hello GoUI",
		Width:     400,
		Height:    300,
		Resizable: true,
	})
	if err != nil {
		panic(err)
	}

	// 2. Set the Root Layout
	// Vertical Box Layout with 20px padding and 10px spacing
	win.Root.SetLayout(&layout.VBoxLayout{
		Padding: 20,
		Spacing: 10,
	})

	// 3. Add Components
	// Create a Label
	label := component.NewLabel("Welcome to GoUI!")
	win.Add(label)

	// Create a Button
	btn := component.NewButton("Click Me")
	
	// Add Click Handler
	// Note: In a real app, you might use a more robust event system
	// but currently Button doesn't expose a direct OnClick field in the struct
	// You typically handle events by checking the EventBus or subclassing.
	// For simplicity in this version, we just add it to the UI.
	
	win.Add(btn)

	// 4. Show and Run
	win.Show()
	win.Run()
}
```

## Running the Examples

You can try the included examples to see GoUI in action:

```bash
# Todo App
go run examples/todo/main.go

# Notepad
go run examples/notepad/main.go

# System Monitor
go run examples/sysmon/main.go

# Snake Game
go run examples/snake/main.go
```
