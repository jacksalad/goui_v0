# Event Handling

GoUI's event system bridges the gap between raw Windows messages and your Go application logic.

## 🔄 The Event Loop

1.  **Windows Message**: The OS sends a message (e.g., `WM_LBUTTONDOWN`) to the window handle.
2.  **WindowProc**: The `window.go` implementation receives this.
3.  **Event Conversion**: It's converted to a `goui/event.Event` struct.
4.  **Distribution**: The event is routed to the appropriate component.

## 📡 Event Types

All events share a common structure:

```go
type Event struct {
    Type      EventType
    Timestamp int64
    Data      interface{} // Type-specific data
}
```

| Event Type | Data Type | Description |
| :--- | :--- | :--- |
| `EventMouseMove` | `MouseEvent` | Mouse moved. `X`, `Y` coords relative to window. |
| `EventMouseClick` | `MouseEvent` | Mouse button pressed. |
| `EventMouseRelease` | `MouseEvent` | Mouse button released. |
| `EventMouseWheel` | `MouseEvent` | Scroll wheel turned. Check `Delta`. |
| `EventKeyPress` | `KeyEvent` | Key pressed. `VirtualKeyCode` (e.g., VK_RETURN). |
| `EventKeyRelease` | `KeyEvent` | Key released. |
| `EventChar` | `KeyEvent` | Character typed. `Rune` contains the char. |
| `EventResize` | `nil` | Window was resized. |
| `EventClose` | `nil` | Window is closing. |

## 🎯 Handling Events in Components

To make a component interactive, override `OnEvent`.

**Return Value:**
*   `true`: "I have handled this event. Stop processing."
*   `false`: "I ignored this event. Let others handle it / Default behavior."

**Example: A Custom Button**

```go
func (b *MyButton) OnEvent(evt event.Event) bool {
    // 1. Check if Mouse Click
    if evt.Type == event.EventMouseClick {
        // 2. Hit Test (Check if click is inside my bounds)
        // Note: The Window usually only dispatches to us if we are under the cursor,
        // or if we have captured focus.
        
        // Check bounds logic is usually handled by the Window dispatcher 
        // before calling OnEvent for mouse events.
        // But for safety/custom logic:
        if b.Bounds.Contains(evt.Data.(MouseEvent).X, evt.Data.(MouseEvent).Y) {
            fmt.Println("Button Clicked!")
            return true
        }
    }
    return false
}
```

## 🎹 Keyboard Focus

Keyboard events (`KeyPress`, `Char`) are **only** sent to the component that currently has **Focus**.

*   **Setting Focus**: Call `window.SetFocus(component)`.
*   **Click-to-Focus**: The default `Window` logic automatically sets focus to a component when it is clicked.
*   **Focus Visuals**: Components should override `OnFocus()` and `OnBlur()` to update their visual state (e.g., draw a border, show a cursor).

```go
func (t *TextBox) OnFocus() {
    t.HasFocus = true
    t.RequestRepaint() // Redraw to show cursor
}

func (t *TextBox) OnBlur() {
    t.HasFocus = false
    t.RequestRepaint() // Redraw to hide cursor
}
```

## 🌐 The Global Event Bus

Sometimes you want to listen to events globally (e.g., global hotkeys, logging). The `Window` exposes an `EventBus`.

```go
// Subscribe
sub := win.EventBus.Subscribe(event.EventKeyPress)

go func() {
    for evt := range sub {
        // Handle global key press
    }
}()
```

**Note**: The EventBus receives a copy of the event *after* it has been processed by the UI components.
