# GoUI Documentation

Welcome to the official documentation for **GoUI**, a lightweight, pure Go GUI library for Windows applications.

## � Introduction

GoUI is designed to be a simple, dependency-free alternative to heavy UI frameworks like Electron or Qt for small to medium-sized Windows tools. It leverages the native Windows API (`user32.dll`, `gdi32.dll`, `kernel32.dll`) directly through Go's `syscall` package, meaning:

*   **No CGo Required**: Pure Go compilation, easy cross-compilation (in theory, though runtime is Windows-only).
*   **Tiny Binary Size**: No bundled browser engine or large runtime.
*   **Native Performance**: Direct GDI drawing and message handling.

## 🌟 Key Features

*   **Retained Mode UI**: Components are objects that persist state, but rendering is immediate-mode style (repaint on dirty).
*   **Flexbox Layout**: A powerful layout engine inspired by CSS Flexbox for responsive designs.
*   **Thread-Safe Concurrency**: Built-in `Window.RequestRepaint()` for safe UI updates from goroutines.
*   **Customizable Components**: Easy to extend `BaseComponent` to create custom widgets.
*   **Rich Standard Library**: Includes Charts, ProgressBars, Cards, and more out of the box.

## 📚 Documentation Sections

1.  [**Getting Started**](getting_started.md)
    *   Installation & Prerequisites
    *   Your First Application
    *   Running Examples
2.  [**Core Concepts & Architecture**](core_concepts.md)
    *   The Window Loop
    *   Component Lifecycle
    *   Rendering Pipeline
    *   Thread Safety & Concurrency
3.  [**Components Reference**](components.md)
    *   **Basic**: Button, Label, Image, CheckBox
    *   **Input**: TextBox, TextArea
    *   **Containers**: Panel, Card
    *   **Data**: ProgressBar, LineChart
4.  [**Layout System**](layouts.md)
    *   Box Model (VBox/HBox)
    *   Flex Layout (Deep Dive)
    *   Grid Layout
5.  [**Event Handling**](events.md)
    *   Event Propagation
    *   Keyboard & Mouse Events
    *   Global Event Bus

## 📦 Project Structure

A typical GoUI application structure:

```
my-app/
├── main.go           # Entry point
├── ui/               # Custom UI components
│   └── my_widget.go
├── assets/           # Images, fonts, etc.
└── go.mod
```

## ⚠️ Status

GoUI is currently in **Alpha**. The API may change, and it is primarily focused on Windows desktop development.

## 🤝 Contributing

We welcome contributions! Whether it's fixing bugs, adding new components, or improving documentation.
