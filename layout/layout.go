package layout

// Interface for layout management
type Layout interface {
	Arrange(container Container)
}

// Container interface to avoid circular dependency with component package
type Container interface {
	GetBounds() Rect
	GetChildren() []Component
}

type Rect struct {
	X, Y, Width, Height int32
}

func (r Rect) Contains(x, y int32) bool {
	return x >= r.X && x < r.X+r.Width && y >= r.Y && y < r.Y+r.Height
}

type Component interface {
	SetBounds(x, y, width, height int32)
	GetBounds() Rect
	GetPreferredSize() (int32, int32)
	IsVisible() bool
}

// VBoxLayout arranges components vertically
type VBoxLayout struct {
	Spacing int32
	Padding int32
}

func (l *VBoxLayout) Arrange(container Container) {
	bounds := container.GetBounds()
	x := bounds.X + l.Padding
	y := bounds.Y + l.Padding
	width := bounds.Width - 2*l.Padding

	for _, child := range container.GetChildren() {
		if !child.IsVisible() {
			continue
		}

		// Use preferred height
		_, prefH := child.GetPreferredSize()
		// If preferred height is 0, use current height (simplistic)
		if prefH <= 0 {
			prefH = child.GetBounds().Height
		}

		child.SetBounds(x, y, width, prefH)
		y += prefH + int32(l.Spacing)
	}
}

// HBoxLayout arranges components horizontally
type HBoxLayout struct {
	Spacing int32
	Padding int32
}

func (l *HBoxLayout) Arrange(container Container) {
	bounds := container.GetBounds()
	x := bounds.X + l.Padding
	y := bounds.Y + l.Padding
	height := bounds.Height - 2*l.Padding

	for _, child := range container.GetChildren() {
		if !child.IsVisible() {
			continue
		}

		prefW, _ := child.GetPreferredSize()
		if prefW <= 0 {
			prefW = child.GetBounds().Width
		}

		child.SetBounds(x, y, prefW, height)
		x += prefW + int32(l.Spacing)
	}
}
