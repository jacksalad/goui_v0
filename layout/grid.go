package layout

type GridProps struct {
	Row, Col         int
	RowSpan, ColSpan int
}

type GridLayout struct {
	Rows, Cols int
	Spacing    int32
	Padding    int32

	// Store child positions
	props map[Component]GridProps
}

func NewGridLayout(rows, cols int) *GridLayout {
	return &GridLayout{
		Rows:  rows,
		Cols:  cols,
		props: make(map[Component]GridProps),
	}
}

func (l *GridLayout) SetPosition(c Component, row, col int) {
	if l.props == nil {
		l.props = make(map[Component]GridProps)
	}
	p := l.props[c]
	p.Row = row
	p.Col = col
	l.props[c] = p
}

func (l *GridLayout) SetSpan(c Component, rowSpan, colSpan int) {
	if l.props == nil {
		l.props = make(map[Component]GridProps)
	}
	p := l.props[c]
	p.RowSpan = rowSpan
	p.ColSpan = colSpan
	l.props[c] = p
}

func (l *GridLayout) Arrange(container Container) {
	bounds := container.GetBounds()
	width := bounds.Width - 2*l.Padding
	height := bounds.Height - 2*l.Padding
	startX := bounds.X + l.Padding
	startY := bounds.Y + l.Padding

	if l.Rows == 0 || l.Cols == 0 {
		return
	}

	// Calculate cell size (Uniform)
	// Avoid divide by zero
	if l.Cols == 0 { l.Cols = 1 }
	if l.Rows == 0 { l.Rows = 1 }
	
	cellW := (width - int32(l.Cols-1)*l.Spacing) / int32(l.Cols)
	cellH := (height - int32(l.Rows-1)*l.Spacing) / int32(l.Rows)

	if cellW < 0 {
		cellW = 0
	}
	if cellH < 0 {
		cellH = 0
	}

	for _, child := range container.GetChildren() {
		if !child.IsVisible() {
			continue
		}

		p, ok := l.props[child]
		if !ok {
			// If not positioned, maybe auto-flow?
			// For now, skip or put at 0,0
			p = GridProps{Row: 0, Col: 0, RowSpan: 1, ColSpan: 1}
		}
		if p.RowSpan == 0 {
			p.RowSpan = 1
		}
		if p.ColSpan == 0 {
			p.ColSpan = 1
		}

		// Calculate position
		x := startX + int32(p.Col)*(cellW+l.Spacing)
		y := startY + int32(p.Row)*(cellH+l.Spacing)

		// Calculate size
		w := int32(p.ColSpan)*cellW + int32(p.ColSpan-1)*l.Spacing
		h := int32(p.RowSpan)*cellH + int32(p.RowSpan-1)*l.Spacing

		child.SetBounds(x, y, w, h)
	}
}
