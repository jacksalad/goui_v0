# Components Reference

This section details the available components, their properties, and usage patterns.

## 🏗️ Base Components

### Button

A standard push button.

**Key Properties:**
*   `Text` (string): The label on the button.
*   `OnClick` (func): *Recommended pattern is to handle `EventMouseClick` in `OnEvent` or use a wrapper.*

**Usage:**
```go
btn := component.NewButton("Save")
// Custom event handling
// (See Event Handling section for advanced patterns)
```

### Label

Displays read-only text.

**Key Properties:**
*   `Text` (string): The content.
*   `Color` (uint32): Text color (0xAARRGGBB format).
*   `Font` (*render.Font): Custom font.

**Usage:**
```go
lbl := component.NewLabel("Status: Ready")
lbl.Color = 0xFF008000 // Green
```

### Image

Displays a bitmap image.

**Supported Formats:** JPG, PNG.

**Usage:**
```go
img, err := component.NewImage("icon.png")
if err != nil { panic(err) }
// Images automatically size to their content unless bounds are set manually
```

## 📝 Input Components

### TextBox

Single-line text input field.

**Key Properties:**
*   `Text` (string): Current value.
*   `BgColor` (uint32): Background color.
*   `TextColor` (uint32): Text color.

**Events Handled:**
*   `EventKeyPress`: Handles Backspace, Delete, Left/Right arrows.
*   `EventChar`: Appends characters.
*   `EventMouseClick`: Sets focus.

**Usage:**
```go
input := component.NewTextBox(200, 30)
input.Text = "Default Value"
```

### TextArea

Multi-line text editor with scroll capabilities.

**Key Properties:**
*   `Text` (string): Full content.
*   `FontSize` (int): Font size.

**Features:**
*   **Scrolling**: Supports mouse wheel and arrow keys.
*   **Cursor**: Blinking cursor tracking position.
*   **Selection**: (Planned/Partial) Text selection support.

**Usage:**
```go
editor := component.NewTextArea(400, 300)
editor.SetText("Line 1\nLine 2")
```

### CheckBox

A toggle widget.

**Key Properties:**
*   `Checked` (bool): State.
*   `Text` (string): Label text.

**Usage:**
```go
chk := component.NewCheckBox("Enable Logging")
if chk.Checked {
    // ...
}
```

## 📦 Containers

### Panel

The fundamental container. Can nest other components.

**Key Properties:**
*   `BgColor` (uint32): Background color.
*   `Layout` (layout.Layout): The layout manager.

**Usage:**
```go
panel := component.NewPanel(0, 0, 100, 100)
panel.SetLayout(&layout.VBoxLayout{})
panel.Add(child1)
panel.Add(child2)
```

### Card

A styled Panel with a title bar, border, and shadow.

**Usage:**
```go
card := component.NewCard(300, 200, "User Profile")
card.Add(component.NewLabel("Name: John Doe"))
```

## 📊 Data Visualization

### ProgressBar

Visualizes a percentage value.

**Key Properties:**
*   `Value` (float64): 0.0 to 1.0.
*   `Color` (uint32): Bar color.

**Usage:**
```go
bar := component.NewProgressBar(200, 20)
bar.SetValue(0.5) // 50%
```

### LineChart

Draws a historical trend line.

**Key Properties:**
*   `Data` ([]float64): The data points.
*   `Color` (uint32): Line color.
*   `MaxPoints` (int): Window size for scrolling data.
*   `MinY`, `MaxY` (float64): Y-axis range.

**Usage:**
```go
chart := component.NewLineChart(300, 150)
chart.Color = 0xFF0000FF
chart.MaxY = 100.0
chart.AddPoint(42.0)
```

## 🛠️ Creating Custom Components

To create a custom component, embed `BaseComponent` and implement `Render` and `OnEvent`.

```go
type MyWidget struct {
    component.BaseComponent
    ClickCount int
}

func NewMyWidget() *MyWidget {
    w := &MyWidget{}
    w.Visible = true
    return w
}

func (w *MyWidget) Render(canvas *render.Canvas) {
    // 1. Draw Background
    canvas.FillRect(w.Bounds, 0xFFEEEEEE)
    
    // 2. Draw Text
    text := fmt.Sprintf("Clicks: %d", w.ClickCount)
    canvas.DrawText(w.Bounds.X + 10, w.Bounds.Y + 10, text, nil, 0xFF000000)
}

func (w *MyWidget) OnEvent(evt event.Event) bool {
    if evt.Type == event.EventMouseClick {
        w.ClickCount++
        w.RequestRepaint() // Important!
        return true
    }
    return false
}
```
