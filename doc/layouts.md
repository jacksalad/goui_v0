# Layout System

GoUI provides a flexible layout system to arrange components automatically. This avoids the need for hardcoded pixel coordinates, making your UI responsive to window resizing.

## 📏 The Layout Interface

Any struct implementing this interface can be a layout manager:

```go
type Layout interface {
    Arrange(container Container)
}
```

The `Container` interface provides access to the children and the container's own bounds:

```go
type Container interface {
    GetBounds() Rect
    GetChildren() []Component
}
```

## 📦 Standard Layouts

### 1. VBoxLayout & HBoxLayout

These are simple stack layouts.

*   **VBox**: Stacks children vertically.
*   **HBox**: Stacks children horizontally.

**Properties:**
*   `Spacing` (int32): Pixels between items.
*   `Padding` (int32): Pixels around the edge of the container.

**Example:**
```go
vbox := &layout.VBoxLayout{
    Spacing: 10,
    Padding: 20,
}
panel.SetLayout(vbox)
```

### 2. FlexLayout

The most powerful layout, inspired by CSS Flexbox. It allows complex weighted sizing and alignment.

**Concepts:**
*   **Main Axis**: The primary direction (Row or Column).
*   **Cross Axis**: The perpendicular direction.
*   **Grow**: How much extra space a child should take.

**Properties:**
*   `Direction`: `FlexRow` or `FlexColumn`.
*   `JustifyContent`: Distribution along Main Axis (`FlexStart`, `FlexCenter`, `FlexSpaceBetween`, etc.).
*   `AlignItems`: Alignment along Cross Axis (`AlignStart`, `AlignCenter`, `AlignStretch`).

**Using Flex Weights (Grow):**

To make a child expand to fill available space, use `SetGrow`.

```go
flex := layout.NewFlexLayout(layout.FlexRow)

// Sidebar: Fixed width
sidebar := component.NewPanel(0, 0, 200, 0)

// Content: Fills remaining space
content := component.NewPanel(0, 0, 0, 0)
flex.SetGrow(content, 1) // Grow factor 1

panel.SetLayout(flex)
panel.Add(sidebar)
panel.Add(content)
```

### 3. GridLayout

Arranges items in a fixed grid.

**Properties:**
*   `Rows`, `Cols`: Number of rows and columns.
*   `Spacing`: Gap between cells.

**Positioning:**
Currently, `GridLayout` automatically flows items into cells (Row 0, Col 0 -> Row 0, Col 1...). 
Future versions may expose explicit `SetRow/SetCol` APIs more publicly.

## 📐 Preferred Size

Layouts rely on `GetPreferredSize()` to know how big a component *wants* to be.

*   **Buttons/Labels**: Calculate size based on text width + padding.
*   **Images**: Use image dimensions.
*   **Panels**: Usually 0,0 unless they contain children, in which case they might calculate a size based on content.

**Tip:** If a component isn't showing up, check if its preferred size is 0 and the layout isn't stretching it!
