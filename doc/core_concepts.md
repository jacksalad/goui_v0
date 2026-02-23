# Core Concepts & Architecture

Understanding how GoUI works under the hood will help you build more robust applications.

## 1. The Component Model

At the heart of GoUI is the `Component` interface. Every UI element—from a simple Button to a complex Chart—implements this interface.

### The `Component` Interface

```go
type Component interface {
    // Drawing
    Render(canvas *render.Canvas)
    
    // Event Handling
    OnEvent(event event.Event) bool
    
    // Layout & Geometry
    SetBounds(x, y, width, height int32)
    GetBounds() layout.Rect
    GetPreferredSize() (width, height int32)
    
    // Visibility & State
    SetVisible(visible bool)
    IsVisible() bool
    
    // Focus Management
    OnFocus()
    OnBlur()
    
    // Lifecycle
    RequestRepaint()
}
```

### BaseComponent
Most components embed `BaseComponent`, which handles common tasks like storing bounds, visibility state, and repaint requests. When creating a custom component, you typically only need to override `Render` and `OnEvent`.

## 2. The Rendering Pipeline

GoUI uses a **Hybrid Retained/Immediate** rendering model.

1.  **Retained State**: The component tree (Window -> Root Panel -> Children) persists in memory. You modify the state of components (e.g., `label.Text = "New"`) directly.
2.  **Dirty Checking**: When state changes, components call `RequestRepaint()`. This sets a flag or posts a message.
3.  **Immediate Rendering**: When the Windows `WM_PAINT` message is received (or manually triggered via `WM_USER`), the `Window` iterates through the component tree and calls `Render(canvas)` on every visible component.

### Double Buffering
(Note: Implementation specific) To avoid flickering, the `Renderer` typically draws to an off-screen bitmap (buffer) first, and then "blits" the entire frame to the window in one go (`BitBlt`).

## 3. The Main Loop & Events

GoUI relies on the standard Windows Message Loop:

```go
for {
    GetMessage(&msg)      // Blocks until a message arrives
    TranslateMessage(&msg) // Translates virtual keys to characters
    DispatchMessage(&msg)  // Sends message to WindowProc
}
```

### Event Flow
1.  **OS Event**: User clicks mouse.
2.  **WindowProc**: Receives `WM_LBUTTONDOWN`.
3.  **GoUI Conversion**: Converts native `MSG` to GoUI `event.Event`.
4.  **Dispatch**:
    *   **Focus Dispatch**: If a component has focus (e.g., TextBox), keyboard events go there first.
    *   **Hit Testing**: For mouse events, GoUI traverses the component tree to find which component is under the cursor.
    *   **Bubbling**: (If implemented) Events can bubble up from child to parent. currently, GoUI primarily uses direct dispatch.

## 4. Thread Safety & Concurrency

**Rule #1: Never update the UI directly from a background goroutine.**
Windows GDI handles are thread-affine. Calling GDI functions from a different thread than the one that created the window can lead to crashes or undefined behavior.

### The Solution: `RequestRepaint`

GoUI provides a thread-safe mechanism to trigger updates:

```go
// Safe to call from ANY goroutine
window.RequestRepaint()
```

**How it works:**
1.  `RequestRepaint` calls `PostMessageW` with `WM_USER`.
2.  `PostMessageW` is thread-safe and puts a message in the UI thread's queue.
3.  The main loop picks up `WM_USER`.
4.  The main loop calls `window.Render()`, which runs on the UI thread.

### Pattern: Updating Data vs Updating UI

It **is** safe to update Go struct fields from background threads (with proper mutex locking if needed), as long as you don't call drawing functions.

**Correct Pattern:**
```go
// Background Goroutine
go func() {
    newData := fetchFromNetwork()
    
    // Lock if shared state
    myComponent.Lock()
    myComponent.Data = newData
    myComponent.Unlock()
    
    // Schedule a repaint
    window.RequestRepaint()
}()
```

## 5. Layout System

Layouts in GoUI are distinct from Components. A `Panel` (Container) has a `Layout` interface.

```go
type Layout interface {
    Arrange(container Container)
}
```

When `Arrange` is called (usually before Render), the layout manager calculates the position and size (`SetBounds`) of all children based on:
1.  The container's size.
2.  The children's `GetPreferredSize()`.
3.  Layout-specific properties (padding, spacing, flex grow, etc.).
