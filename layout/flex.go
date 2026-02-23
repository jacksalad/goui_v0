package layout

// FlexDirection defines the direction of the flex container
type FlexDirection int

const (
	FlexRow FlexDirection = iota
	FlexColumn
)

type FlexJustify int

const (
	FlexStart FlexJustify = iota
	FlexEnd
	FlexCenter
	FlexSpaceBetween
	FlexSpaceAround
)

type FlexAlign int

const (
	AlignStart FlexAlign = iota
	AlignEnd
	AlignCenter
	AlignStretch
)

type FlexProps struct {
	Grow   int
	Shrink int
}

type FlexLayout struct {
	Direction      FlexDirection
	JustifyContent FlexJustify
	AlignItems     FlexAlign
	Spacing        int32
	Padding        int32

	// Layout properties for children
	props map[Component]FlexProps
}

func NewFlexLayout(direction FlexDirection) *FlexLayout {
	return &FlexLayout{
		Direction: direction,
		props:     make(map[Component]FlexProps),
	}
}

func (l *FlexLayout) SetGrow(c Component, grow int) {
	// We need to store props, but map keys must be comparable.
	// Interfaces are comparable.
	if l.props == nil {
		l.props = make(map[Component]FlexProps)
	}
	p := l.props[c]
	p.Grow = grow
	l.props[c] = p
}

func (l *FlexLayout) Arrange(container Container) {
	bounds := container.GetBounds()
	width := bounds.Width - 2*l.Padding
	height := bounds.Height - 2*l.Padding
	startX := bounds.X + l.Padding
	startY := bounds.Y + l.Padding

	if l.Direction == FlexRow {
		l.arrangeRow(container, startX, startY, width, height)
	} else {
		l.arrangeColumn(container, startX, startY, width, height)
	}
}

func (l *FlexLayout) arrangeRow(container Container, x, y, width, height int32) {
	children := container.GetChildren()
	if len(children) == 0 {
		return
	}

	// 1. Measure total fixed width and total grow
	totalFixedWidth := int32(0)
	totalGrow := 0
	visibleCount := 0

	childSizes := make([]struct{ w, h int32 }, len(children))

	for i, child := range children {
		if !child.IsVisible() {
			continue
		}
		visibleCount++
		w, h := child.GetPreferredSize()
		if w <= 0 {
			w = child.GetBounds().Width
		}
		if h <= 0 {
			h = child.GetBounds().Height
		}

		childSizes[i] = struct{ w, h int32 }{w, h}

		totalFixedWidth += w
		if p, ok := l.props[child]; ok {
			totalGrow += p.Grow
		}
	}

	if visibleCount > 0 {
		totalFixedWidth += int32(visibleCount-1) * l.Spacing
	}

	remainingSpace := width - totalFixedWidth
	if remainingSpace < 0 {
		remainingSpace = 0 // Overflow, we don't handle shrink yet properly, just clip
	}

	// 2. Distribute space
	currentX := x

	// Calculate gap for justification if no grow
	justifyGap := int32(0)
	if totalGrow == 0 && remainingSpace > 0 {
		switch l.JustifyContent {
		case FlexEnd:
			currentX += remainingSpace
		case FlexCenter:
			currentX += remainingSpace / 2
		case FlexSpaceBetween:
			if visibleCount > 1 {
				justifyGap = remainingSpace / int32(visibleCount-1)
			}
		case FlexSpaceAround:
			if visibleCount > 0 {
				justifyGap = remainingSpace / int32(visibleCount) // Approximate
				currentX += justifyGap / 2
			}
		}
	}

	for i, child := range children {
		if !child.IsVisible() {
			continue
		}

		w, h := childSizes[i].w, childSizes[i].h

		// Apply Grow
		grow := 0
		if p, ok := l.props[child]; ok {
			grow = p.Grow
		}

		if totalGrow > 0 && remainingSpace > 0 {
			extra := int32(float32(remainingSpace) * (float32(grow) / float32(totalGrow)))
			w += extra
		}

		// Apply AlignItems (Cross Axis)
		childY := y
		childH := h

		switch l.AlignItems {
		case AlignStart:
			childY = y
		case AlignEnd:
			childY = y + height - h
		case AlignCenter:
			childY = y + (height-h)/2
		case AlignStretch:
			childY = y
			childH = height
		}

		child.SetBounds(currentX, childY, w, childH)

		currentX += w + l.Spacing + justifyGap
	}
}

func (l *FlexLayout) arrangeColumn(container Container, x, y, width, height int32) {
	children := container.GetChildren()
	if len(children) == 0 {
		return
	}

	// 1. Measure
	totalFixedHeight := int32(0)
	totalGrow := 0
	visibleCount := 0

	childSizes := make([]struct{ w, h int32 }, len(children))

	for i, child := range children {
		if !child.IsVisible() {
			continue
		}
		visibleCount++
		w, h := child.GetPreferredSize()
		if w <= 0 {
			w = child.GetBounds().Width
		}
		if h <= 0 {
			h = child.GetBounds().Height
		}

		childSizes[i] = struct{ w, h int32 }{w, h}

		totalFixedHeight += h
		if p, ok := l.props[child]; ok {
			totalGrow += p.Grow
		}
	}

	if visibleCount > 0 {
		totalFixedHeight += int32(visibleCount-1) * l.Spacing
	}

	remainingSpace := height - totalFixedHeight
	if remainingSpace < 0 {
		remainingSpace = 0
	}

	// 2. Distribute space
	currentY := y

	justifyGap := int32(0)
	if totalGrow == 0 && remainingSpace > 0 {
		switch l.JustifyContent {
		case FlexEnd:
			currentY += remainingSpace
		case FlexCenter:
			currentY += remainingSpace / 2
		case FlexSpaceBetween:
			if visibleCount > 1 {
				justifyGap = remainingSpace / int32(visibleCount-1)
			}
		case FlexSpaceAround:
			if visibleCount > 0 {
				justifyGap = remainingSpace / int32(visibleCount)
				currentY += justifyGap / 2
			}
		}
	}

	for i, child := range children {
		if !child.IsVisible() {
			continue
		}

		w, h := childSizes[i].w, childSizes[i].h

		// Apply Grow
		grow := 0
		if p, ok := l.props[child]; ok {
			grow = p.Grow
		}

		if totalGrow > 0 && remainingSpace > 0 {
			extra := int32(float32(remainingSpace) * (float32(grow) / float32(totalGrow)))
			h += extra
		}

		// Apply AlignItems (Cross Axis)
		childX := x
		childW := w

		switch l.AlignItems {
		case AlignStart:
			childX = x
		case AlignEnd:
			childX = x + width - w
		case AlignCenter:
			childX = x + (width-w)/2
		case AlignStretch:
			childX = x
			childW = width
		}

		child.SetBounds(childX, currentY, childW, h)

		currentY += h + l.Spacing + justifyGap
	}
}
