package component

import (
	"github.com/jacksalad/goui_v0/render"
)

type ProgressBar struct {
	BaseComponent
	Value   float64 // 0.0 to 1.0
	Color   uint32
	BgColor uint32
}

func NewProgressBar(width, height int32) *ProgressBar {
	p := &ProgressBar{
		Value:   0.0,
		Color:   0xFF0078D7, // Blue
		BgColor: 0xFFE0E0E0, // Gray
	}
	p.SetBounds(0, 0, width, height)
	p.Visible = true
	return p
}

func (p *ProgressBar) SetValue(v float64) {
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	p.Value = v
	p.RequestRepaint()
}

func (p *ProgressBar) Render(canvas *render.Canvas) {
	if !p.Visible {
		return
	}

	// Background
	canvas.FillRect(p.Bounds.X, p.Bounds.Y, p.Bounds.Width, p.Bounds.Height, p.BgColor)

	// Foreground
	fillWidth := int32(float64(p.Bounds.Width) * p.Value)
	if fillWidth > 0 {
		canvas.FillRect(p.Bounds.X, p.Bounds.Y, fillWidth, p.Bounds.Height, p.Color)
	}

	// Border
	borderColor := uint32(0xFFAAAAAA)
	canvas.FillRect(p.Bounds.X, p.Bounds.Y, p.Bounds.Width, 1, borderColor)
	canvas.FillRect(p.Bounds.X, p.Bounds.Y+p.Bounds.Height-1, p.Bounds.Width, 1, borderColor)
	canvas.FillRect(p.Bounds.X, p.Bounds.Y, 1, p.Bounds.Height, borderColor)
	canvas.FillRect(p.Bounds.X+p.Bounds.Width-1, p.Bounds.Y, 1, p.Bounds.Height, borderColor)
}
