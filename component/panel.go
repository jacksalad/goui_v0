package component

import (
	"github.com/jacksalad/goui_v0/event"
	"github.com/jacksalad/goui_v0/layout"
	"github.com/jacksalad/goui_v0/render"
)

type Panel struct {
	BaseComponent
	Children []Component
	BgColor  uint32
	Layout   layout.Layout
}

func NewPanel(x, y, w, h int32) *Panel {
	p := &Panel{
		Children: make([]Component, 0),
		BgColor:  0xFFFFFFFF, // White
	}
	p.SetBounds(x, y, w, h)
	p.Visible = true
	return p
}

func (p *Panel) GetPreferredSize() (int32, int32) {
	// Simple preferred size based on layout or children?
	// For now, return current bounds
	return p.Bounds.Width, p.Bounds.Height
}

func (p *Panel) GetBounds() layout.Rect {
	return layout.Rect{
		X:      p.Bounds.X,
		Y:      p.Bounds.Y,
		Width:  p.Bounds.Width,
		Height: p.Bounds.Height,
	}
}

func (p *Panel) GetChildren() []layout.Component {
	res := make([]layout.Component, len(p.Children))
	for i, c := range p.Children {
		res[i] = c
	}
	return res
}

func (p *Panel) Add(c Component) {
	p.Children = append(p.Children, c)
	p.LayoutChildren()
}

func (p *Panel) Remove(c Component) {
	for i, child := range p.Children {
		if child == c {
			p.Children = append(p.Children[:i], p.Children[i+1:]...)
			p.LayoutChildren()
			return
		}
	}
}

func (p *Panel) SetLayout(l layout.Layout) {
	p.Layout = l
	p.LayoutChildren()
}

func (p *Panel) LayoutChildren() {
	if p.Layout != nil {
		p.Layout.Arrange(p)
	}
}

// Override SetBounds to re-layout when panel resizes
func (p *Panel) SetBounds(x, y, width, height int32) {
	p.BaseComponent.SetBounds(x, y, width, height)
	p.LayoutChildren()
}

func (p *Panel) Render(canvas *render.Canvas) {
	if !p.Visible {
		return
	}
	// Fill background
	canvas.FillRect(p.Bounds.X, p.Bounds.Y, p.Bounds.Width, p.Bounds.Height, p.BgColor)

	// Render children
	for _, child := range p.Children {
		child.Render(canvas)
	}
	p.RepaintRequested = false
}

func (p *Panel) FindComponentAt(x, y int32) Component {
	if !p.Visible || !p.Bounds.Contains(x, y) {
		return nil
	}

	for i := len(p.Children) - 1; i >= 0; i-- {
		child := p.Children[i]

		// If child is a Panel, recurse
		if panel, ok := child.(*Panel); ok {
			if found := panel.FindComponentAt(x, y); found != nil {
				return found
			}
		} else {
			// Leaf component check
			rect := child.GetBounds()
			if rect.Contains(x, y) {
				return child
			}
		}
	}
	return p
}

func (p *Panel) OnEvent(evt event.Event) bool {
	if !p.Visible {
		return false
	}

	// Dispatch to children in reverse order (top-most first)
	for i := len(p.Children) - 1; i >= 0; i-- {
		if p.Children[i].OnEvent(evt) {
			return true
		}
	}
	return false
}
