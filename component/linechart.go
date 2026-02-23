package component

import (
	"goui/render"
)

type LineChart struct {
	BaseComponent
	Data      []float64
	MaxPoints int
	MinY      float64
	MaxY      float64
	Color     uint32
	BgColor   uint32
}

func NewLineChart(width, height int32) *LineChart {
	l := &LineChart{
		Data:      make([]float64, 0),
		MaxPoints: 100,
		MinY:      0,
		MaxY:      100,
		Color:     0xFF0078D7, // Blue
		BgColor:   0xFFFFFFFF, // White
	}
	l.SetBounds(0, 0, width, height)
	l.Visible = true
	return l
}

func (l *LineChart) AddPoint(v float64) {
	l.Data = append(l.Data, v)
	if len(l.Data) > l.MaxPoints {
		l.Data = l.Data[1:]
	}
	l.RequestRepaint()
}

func (l *LineChart) Render(canvas *render.Canvas) {
	if !l.Visible {
		return
	}

	// Background
	canvas.FillRect(l.Bounds.X, l.Bounds.Y, l.Bounds.Width, l.Bounds.Height, l.BgColor)

	// Border
	borderColor := uint32(0xFFAAAAAA)
	canvas.FillRect(l.Bounds.X, l.Bounds.Y, l.Bounds.Width, 1, borderColor)
	canvas.FillRect(l.Bounds.X, l.Bounds.Y+l.Bounds.Height-1, l.Bounds.Width, 1, borderColor)
	canvas.FillRect(l.Bounds.X, l.Bounds.Y, 1, l.Bounds.Height, borderColor)
	canvas.FillRect(l.Bounds.X+l.Bounds.Width-1, l.Bounds.Y, 1, l.Bounds.Height, borderColor)

	if len(l.Data) < 2 {
		return
	}

	// Calculate scaling
	stepX := float64(l.Bounds.Width) / float64(l.MaxPoints-1)
	rangeY := l.MaxY - l.MinY
	if rangeY <= 0 {
		rangeY = 1
	}
	// scaleY := float64(l.Bounds.Height) / rangeY

	// Draw lines
	for i := 0; i < len(l.Data)-1; i++ {
		x0 := int32(float64(i) * stepX)
		v0 := l.Data[i]
		norm0 := (v0 - l.MinY) / rangeY
		if norm0 < 0 {
			norm0 = 0
		}
		if norm0 > 1 {
			norm0 = 1
		}
		y0 := l.Bounds.Height - int32(norm0*float64(l.Bounds.Height))

		x1 := int32(float64(i+1) * stepX)
		v1 := l.Data[i+1]
		norm1 := (v1 - l.MinY) / rangeY
		if norm1 < 0 {
			norm1 = 0
		}
		if norm1 > 1 {
			norm1 = 1
		}
		y1 := l.Bounds.Height - int32(norm1*float64(l.Bounds.Height))

		canvas.DrawLine(l.Bounds.X+x0, l.Bounds.Y+y0, l.Bounds.X+x1, l.Bounds.Y+y1, l.Color)
	}
}
