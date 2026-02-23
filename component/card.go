package component

import (
	"goui/event"
	"goui/layout"
	"goui/render"
)

type Card struct {
	BaseComponent
	Title      string
	BgColor    uint32
	InnerPanel *Panel
}

func NewCard(width, height int32, title string) *Card {
	c := &Card{
		Title:   title,
		BgColor: 0xFFFFFFFF,
	}
	c.SetBounds(0, 0, width, height)
	c.Visible = true

	// Inner Panel for content
	c.InnerPanel = NewPanel(10, 35, width-20, height-45)
	c.InnerPanel.BgColor = 0xFFFFFFFF // Transparent? Or same as Card

	// Default layout for inner panel
	c.InnerPanel.SetLayout(&layout.VBoxLayout{Spacing: 5})

	return c
}

func (c *Card) Add(comp Component) {
	c.InnerPanel.Add(comp)
}

func (c *Card) SetLayout(l layout.Layout) {
	c.InnerPanel.SetLayout(l)
}

func (c *Card) Render(canvas *render.Canvas) {
	if !c.Visible {
		return
	}

	// Draw Background
	canvas.FillRect(c.Bounds.X, c.Bounds.Y, c.Bounds.Width, c.Bounds.Height, c.BgColor)

	// Draw Border (Soft Gray)
	borderColor := uint32(0xFFDDDDDD)
	canvas.FillRect(c.Bounds.X, c.Bounds.Y, c.Bounds.Width, 1, borderColor)
	canvas.FillRect(c.Bounds.X, c.Bounds.Y+c.Bounds.Height-1, c.Bounds.Width, 1, borderColor)
	canvas.FillRect(c.Bounds.X, c.Bounds.Y, 1, c.Bounds.Height, borderColor)
	canvas.FillRect(c.Bounds.X+c.Bounds.Width-1, c.Bounds.Y, 1, c.Bounds.Height, borderColor)

	// Draw Title
	if c.Title != "" {
		// Draw Title Text
		// Assuming default font is set in Renderer for now
		canvas.DrawText(c.Bounds.X+10, c.Bounds.Y+8, c.Title, 0xFF333333)

		// Separator line
		canvas.FillRect(c.Bounds.X, c.Bounds.Y+30, c.Bounds.Width, 1, 0xFFEEEEEE)
	}

	// Update InnerPanel bounds relative to Card
	c.InnerPanel.SetBounds(c.Bounds.X+10, c.Bounds.Y+35, c.Bounds.Width-20, c.Bounds.Height-45)
	c.InnerPanel.Render(canvas)
}

func (c *Card) OnEvent(evt event.Event) bool {
	// Pass events to InnerPanel
	return c.InnerPanel.OnEvent(evt)
}
