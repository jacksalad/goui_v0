package component

import (
	"github.com/jacksalad/goui_v0/event"
	"github.com/jacksalad/goui_v0/render"
)

type Button struct {
	BaseComponent
	Text    string
	OnClick func()
	Font    *render.Font

	// State
	isHovered bool
	isPressed bool
}

func NewButton(text string) *Button {
	b := &Button{
		Text: text,
	}
	// Default size based on text + padding
	w, h := b.GetPreferredSize()
	b.SetBounds(0, 0, w, h)
	b.Visible = true
	return b
}

func (b *Button) GetPreferredSize() (int32, int32) {
	w, h := render.MeasureText(b.Text, b.Font)
	// Add padding
	return w + 20, h + 10
}

func (b *Button) Render(canvas *render.Canvas) {
	if !b.Visible {
		return
	}

	bgColor := uint32(0xFFDDDDDD) // Light gray
	if b.isPressed {
		bgColor = 0xFFAAAAAA // Darker gray
	} else if b.isHovered {
		bgColor = 0xFFEEEEEE // Lighter gray
	}

	canvas.FillRect(b.Bounds.X, b.Bounds.Y, b.Bounds.Width, b.Bounds.Height, bgColor)

	// Draw text centered
	if b.Font != nil {
		canvas.SetFont(b.Font)
	}

	textW, textH := render.MeasureText(b.Text, b.Font)

	textX := b.Bounds.X + (b.Bounds.Width-textW)/2
	textY := b.Bounds.Y + (b.Bounds.Height-textH)/2

	canvas.DrawText(textX, textY, b.Text, 0xFF000000)

	b.RepaintRequested = false
}

func (b *Button) OnEvent(evt event.Event) bool {
	if !b.Visible {
		return false
	}

	switch evt.Type {
	case event.EventMouseMove:
		if data, ok := evt.Data.(event.MouseEvent); ok {
			wasHovered := b.isHovered
			b.isHovered = b.Bounds.Contains(data.X, data.Y)
			if wasHovered != b.isHovered {
				return true // Need repaint
			}
		}

	case event.EventMouseClick: // Mouse Down
		if data, ok := evt.Data.(event.MouseEvent); ok {
			if b.Bounds.Contains(data.X, data.Y) {
				b.isPressed = true
				return true
			}
		}

	case event.EventMouseRelease: // Mouse Up
		if data, ok := evt.Data.(event.MouseEvent); ok {
			wasPressed := b.isPressed
			b.isPressed = false
			if wasPressed && b.Bounds.Contains(data.X, data.Y) {
				if b.OnClick != nil {
					b.OnClick()
				}
				return true
			}
			if wasPressed {
				return true // State changed
			}
		}
	}
	return false
}
