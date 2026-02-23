package component

import (
	"goui/event"
	"goui/render"
)

type CheckBox struct {
	BaseComponent
	Text      string
	Checked   bool
	OnCheck   func(checked bool)
	Font      *render.Font
	isHovered bool
}

func NewCheckBox(text string) *CheckBox {
	c := &CheckBox{
		Text: text,
	}
	c.Visible = true
	return c
}

func (c *CheckBox) GetPreferredSize() (int32, int32) {
	w, h := render.MeasureText(c.Text, c.Font)
	boxSize := int32(16)
	if h < boxSize {
		h = boxSize
	}
	return w + boxSize + 8, h
}

func (c *CheckBox) Render(canvas *render.Canvas) {
	if !c.Visible {
		return
	}

	// Draw Box
	boxSize := int32(16)
	boxX := c.Bounds.X
	boxY := c.Bounds.Y + (c.Bounds.Height-boxSize)/2

	// Box background
	bgColor := uint32(0xFFFFFFFF)
	if c.isHovered {
		bgColor = 0xFFEEEEEE
	}
	canvas.FillRect(boxX, boxY, boxSize, boxSize, bgColor)

	// Box Border (simulated with inner rect for now, or just fill slightly smaller)
	// TODO: Add proper DrawRect (stroke) to Canvas
	// For now, let's just fill a smaller white rect inside a gray rect if we want border
	// But let's keep it simple: just a filled rect.

	// Draw Checkmark if checked
	if c.Checked {
		// Simple X or box inside
		innerSize := int32(10)
		innerX := boxX + (boxSize-innerSize)/2
		innerY := boxY + (boxSize-innerSize)/2
		canvas.FillRect(innerX, innerY, innerSize, innerSize, 0xFF000000)
	}

	// Draw Text
	if c.Font != nil {
		canvas.SetFont(c.Font)
	}
	textX := boxX + boxSize + 8
	textY := c.Bounds.Y + (c.Bounds.Height-16)/2
	// Centering text vertically relative to box might be tricky if font size varies
	// Better to use MeasureText to center
	_, textH := render.MeasureText(c.Text, c.Font)
	textY = c.Bounds.Y + (c.Bounds.Height-textH)/2

	canvas.DrawText(textX, textY, c.Text, 0xFF000000)
	c.RepaintRequested = false
}

func (c *CheckBox) OnEvent(evt event.Event) bool {
	if !c.Visible {
		return false
	}

	switch evt.Type {
	case event.EventMouseMove:
		if data, ok := evt.Data.(event.MouseEvent); ok {
			wasHovered := c.isHovered
			c.isHovered = c.Bounds.Contains(data.X, data.Y)
			if wasHovered != c.isHovered {
				return true
			}
		}

	case event.EventMouseClick: // Mouse Down
		if data, ok := evt.Data.(event.MouseEvent); ok {
			if c.Bounds.Contains(data.X, data.Y) {
				c.Checked = !c.Checked
				if c.OnCheck != nil {
					c.OnCheck(c.Checked)
				}
				return true
			}
		}
	}
	return false
}
