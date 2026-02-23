package component

import (
	"github.com/jacksalad/goui_v0/event"
	"github.com/jacksalad/goui_v0/layout"
	"github.com/jacksalad/goui_v0/render"
)

type Component interface {
	Render(canvas *render.Canvas)
	OnEvent(event event.Event) bool
	SetBounds(x, y, width, height int32)
	GetBounds() layout.Rect
	SetVisible(visible bool)
	IsVisible() bool
	RequestRepaint()
	OnFocus()
	OnBlur()
	GetPreferredSize() (int32, int32)
}

type BaseComponent struct {
	Bounds           layout.Rect
	Visible          bool
	RepaintRequested bool // Simplistic dirty flag
}

func (b *BaseComponent) OnFocus() {}
func (b *BaseComponent) OnBlur() {}

func (b *BaseComponent) GetPreferredSize() (int32, int32) {
	return b.Bounds.Width, b.Bounds.Height
}

func (b *BaseComponent) RequestRepaint() {
	b.RepaintRequested = true
}

func (b *BaseComponent) SetBounds(x, y, width, height int32) {
	b.Bounds = layout.Rect{X: x, Y: y, Width: width, Height: height}
}

func (b *BaseComponent) GetBounds() layout.Rect {
	return b.Bounds
}

func (b *BaseComponent) SetVisible(visible bool) {
	b.Visible = visible
}

func (b *BaseComponent) IsVisible() bool {
	return b.Visible
}

func (b *BaseComponent) Render(canvas *render.Canvas)   {}
func (b *BaseComponent) OnEvent(event event.Event) bool { return false }
