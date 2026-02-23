# GoUI

> **A Lightweight, Pure Go GUI Library for Windows**

GoUI is a simple, dependency-free GUI framework for Windows applications, written entirely in Go. It interacts directly with the Windows API (`user32.dll`, `gdi32.dll`) to create native windows and handle events, while providing a custom, high-level component system and rendering engine.

It is designed for educational purposes, small tools, and developers who want to build Windows desktop applications without the overhead of Electron, Qt, or CGo.

## 🚀 Key Features

*   **Pure Go**: No CGo required. Compiles to a single, static binary.
*   **Native Performance**: Direct Win32 API integration for window management and GDI drawing.
*   **Modern Layouts**: Includes a powerful **Flexbox** layout engine, as well as standard **Grid**, **VBox**, and **HBox** layouts.
*   **Rich Components**:
    *   **Basic**: Button, Label, Image, CheckBox, ProgressBar.
    *   **Input**: TextBox, TextArea (multiline).
    *   **Containers**: Panel, Card, ScrollView.
    *   **Data Visualization**: LineChart.
*   **Thread-Safe**: Built-in concurrency support for safe UI updates from background goroutines (`Window.RequestRepaint`).
*   **Customizable**: Easy-to-extend component architecture.

## 📦 Installation

Since GoUI is a library, you can import it into your Go project.

```bash
go get github.com/jacksalad/goui_v0
```

*Note: GoUI currently supports **Windows** only.*

## ⚡ Quick Start

Here is a minimal "Hello World" application:

```go
package main

import (
	"github.com/jacksalad/goui_v0/component"
	"github.com/jacksalad/goui_v0/layout"
	"github.com/jacksalad/goui_v0/window"
)

func main() {
	// 1. Create the main window
	win, err := window.NewWindow(window.WindowConfig{
		Title:  "Hello GoUI",
		Width:  400,
		Height: 300,
	})
	if err != nil {
		panic(err)
	}

	// 2. Set the root layout (Vertical Box)
	win.Root.SetLayout(&layout.VBoxLayout{
		Padding: 20,
		Spacing: 10,
	})

	// 3. Add components
	win.Add(component.NewLabel("Welcome to GoUI!"))
	win.Add(component.NewButton("Click Me"))

	// 4. Run the application loop
	win.Show()
	win.Run()
}
```

## 🎮 Examples

GoUI comes with several examples to demonstrate its capabilities. You can run them directly from the repository:

### 1. Snake Game 🐍
A classic Snake game implementation demonstrating the game loop, custom rendering, and keyboard input.
```bash
go run examples/snake/main.go
```

### 2. Todo App ✅
A simple task manager showing list manipulation, input handling, and layout management.
```bash
go run examples/todo/main.go
```

### 3. System Monitor 📊
Real-time CPU and Memory usage visualization using the `LineChart` component.
```bash
go run examples/sysmon/main.go
```

### 4. Notepad 📝
A basic text editor demonstrating the `TextArea` component and file I/O operations.
```bash
go run examples/notepad/main.go
```

### More Examples
*   **Calculator**: `examples/calculator`
*   **Tic-Tac-Toe**: `examples/tictactoe`
*   **Layout Demo**: `examples/layout`
*   **Components Showcase**: `examples/components`

## 📚 Documentation

Detailed documentation is available in the `doc/` directory:

*   [**Getting Started**](doc/getting_started.md)
*   [**Core Concepts**](doc/core_concepts.md) (Architecture, Rendering, Thread Safety)
*   [**Components Reference**](doc/components.md)
*   [**Layout System**](doc/layouts.md)
*   [**Event Handling**](doc/events.md)

## 📂 Project Structure

```
goui/
├── component/    # UI Widgets (Button, Label, etc.)
├── doc/          # Documentation files
├── event/        # Event definitions and EventBus
├── examples/     # Demo applications
├── layout/       # Layout managers (Flex, Grid, VBox)
├── render/       # Rendering engine (Canvas, GDI wrappers)
├── window/       # Win32 Window creation and Message Loop
└── main.go       # (Optional) Library entry point
```

## 🤝 Contributing

Contributions are welcome! Feel free to submit issues or pull requests to improve components, add new features, or fix bugs.

## 📄 License

This project is licensed under the **MIT License**.
