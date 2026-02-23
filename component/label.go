package component

import (
	"github.com/jacksalad/goui_v0/render"
)

type Label struct {
	BaseComponent
	Text    string
	FgColor uint32
	Font    *render.Font
}

func NewLabel(text string) *Label {
	l := &Label{
		Text:    text,
		FgColor: 0xFF000000, // Black
	}
	l.Visible = true
	// Default size based on text
	w, h := l.GetPreferredSize()
	l.SetBounds(0, 0, w, h)
	return l
}

func (l *Label) Render(canvas *render.Canvas) {
	if !l.Visible {
		return
	}
	// If label has specific font, set it
	if l.Font != nil {
		canvas.SetFont(l.Font)
	}
	canvas.DrawText(l.Bounds.X, l.Bounds.Y, l.Text, l.FgColor)
	l.RepaintRequested = false
}

func (l *Label) GetPreferredSize() (int32, int32) {
	w, h := render.MeasureText(l.Text, l.Font)
	return w, h
}
