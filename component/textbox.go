package component

import (
	"goui/event"
	"goui/render"
	"time"
)

type TextBox struct {
	BaseComponent
	Text        string
	Placeholder string
	Font        *render.Font
	ReadOnly    bool

	// State
	isFocused     bool
	cursorPos     int // Cursor position (index in runes)
	selStart      int // Selection start index (Anchor)
	isDragging    bool
	pendingMouseX int32 // For deferring calculation
	cursorBlink   bool
	lastBlink     int64
}

func NewTextBox(width int32) *TextBox {
	t := &TextBox{
		cursorPos: 0,
		selStart:  0,
	}
	t.SetBounds(0, 0, width, 24) // Default height
	t.Visible = true
	return t
}

func (t *TextBox) GetPreferredSize() (int32, int32) {
	// Measure height based on font
	_, h := render.MeasureText("Tg", t.Font)
	// Add padding (20px top/bottom for larger feel)
	return t.Bounds.Width, h + 10
}

// getIndexFromX returns the character index for a given X coordinate relative to the text box
func (t *TextBox) getIndexFromX(x int32, canvas *render.Canvas) int {
	textX := t.Bounds.X + 10 // Match left padding
	localX := x - textX
	if localX <= 0 {
		return 0
	}

	runes := []rune(t.Text)
	if len(runes) == 0 {
		return 0
	}

	// Iterate to find the best split point
	// This is a bit inefficient but works for short text
	for i := 0; i <= len(runes); i++ {
		w, _ := canvas.MeasureText(string(runes[:i]))
		if w > localX {
			// Check if we are closer to i-1 or i
			prevW, _ := canvas.MeasureText(string(runes[:i-1]))
			if localX-prevW < w-localX {
				return i - 1
			}
			return i
		}
	}
	return len(runes)
}

func (t *TextBox) Render(canvas *render.Canvas) {
	if !t.Visible {
		return
	}

	// Set Font
	if t.Font != nil {
		canvas.SetFont(t.Font)
	}

	// Handle pending mouse interaction (deferred because we need canvas for measurement)
	if t.isDragging || t.pendingMouseX != 0 {
		// Only if dragging or pending click
		// We use pendingMouseX as a signal.
		// Wait, pendingMouseX is set on Click.
		// If dragging, we update pendingMouseX.
		if t.isDragging {
			idx := t.getIndexFromX(t.pendingMouseX, canvas)
			t.cursorPos = idx
			// If starting drag (selStart == -1), set anchor
			if t.selStart == -1 {
				t.selStart = idx
			}
		} else if t.pendingMouseX != 0 && t.selStart == -1 {
			// Just a click (selStart == -1 from OnEvent)
			idx := t.getIndexFromX(t.pendingMouseX, canvas)
			t.cursorPos = idx
			t.selStart = idx
			t.pendingMouseX = 0 // Clear pending
		}
	}

	// Draw Background
	bgColor := uint32(0xFFFFFFFF)
	if t.isFocused {
		bgColor = 0xFFFFFFFF
	} else {
		bgColor = 0xFFF8F8F8
	}
	canvas.FillRect(t.Bounds.X, t.Bounds.Y, t.Bounds.Width, t.Bounds.Height, bgColor)

	// Draw Border
	borderColor := uint32(0xFFAAAAAA)
	if t.isFocused {
		borderColor = 0xFF0078D7 // Blue
	}
	// Simple border: stroke rect
	// Top
	canvas.FillRect(t.Bounds.X, t.Bounds.Y, t.Bounds.Width, 1, borderColor)
	// Bottom
	canvas.FillRect(t.Bounds.X, t.Bounds.Y+t.Bounds.Height-1, t.Bounds.Width, 1, borderColor)
	// Left
	canvas.FillRect(t.Bounds.X, t.Bounds.Y, 1, t.Bounds.Height, borderColor)
	// Right
	canvas.FillRect(t.Bounds.X+t.Bounds.Width-1, t.Bounds.Y, 1, t.Bounds.Height, borderColor)

	// Draw Text
	textX := t.Bounds.X + 10 // Increased left padding

	// Measure text height to center vertically
	_, textH := canvas.MeasureText("Tg")
	textY := t.Bounds.Y + (t.Bounds.Height-textH)/2

	displayText := t.Text
	textColor := uint32(0xFF000000)

	if len(t.Text) == 0 && len(t.Placeholder) > 0 {
		displayText = t.Placeholder
		textColor = 0xFF888888
	}

	// Draw Selection Highlight
	if t.isFocused && t.selStart != t.cursorPos && t.selStart != -1 {
		start := t.selStart
		end := t.cursorPos
		if start > end {
			start, end = end, start
		}

		runes := []rune(t.Text)
		if start < 0 {
			start = 0
		}
		if end > len(runes) {
			end = len(runes)
		}

		startX, _ := canvas.MeasureText(string(runes[:start]))
		endX, _ := canvas.MeasureText(string(runes[:end]))

		selRectX := textX + startX
		selRectW := endX - startX
		canvas.FillRect(selRectX, textY, selRectW, textH, 0xFFADD8E6) // Light Blue
	}

	canvas.DrawText(textX, textY, displayText, textColor)

	// Draw Cursor
	if t.isFocused {
		now := time.Now().UnixNano()
		if now-t.lastBlink > 500*1e6 { // 500ms blink
			t.cursorBlink = !t.cursorBlink
			t.lastBlink = now
			t.RequestRepaint() // Request repaint when cursor blinks
		}

		if t.cursorBlink {
			// Calculate cursor X position
			runes := []rune(t.Text)
			pos := t.cursorPos
			if pos > len(runes) {
				pos = len(runes)
			}
			if pos < 0 {
				pos = 0
			}

			w, _ := canvas.MeasureText(string(runes[:pos]))
			cursorX := textX + w
			canvas.FillRect(cursorX, textY, 2, textH, 0xFF000000)
		}
	}
	t.RepaintRequested = false
}

func (t *TextBox) OnFocus() {
	t.isFocused = true
	t.lastBlink = 0 // Reset blink timer
	t.cursorBlink = true
	t.RequestRepaint()
}

func (t *TextBox) OnBlur() {
	t.isFocused = false
	t.cursorBlink = false
	t.RequestRepaint()
}

func (t *TextBox) OnEvent(evt event.Event) bool {
	if !t.Visible {
		return false
	}

	switch evt.Type {
	case event.EventMouseClick:
		if data, ok := evt.Data.(event.MouseEvent); ok {
			if t.Bounds.Contains(data.X, data.Y) {
				t.pendingMouseX = data.X
				t.isDragging = true // Start potential drag
				t.selStart = -1     // Signal to recalculate in Render (start anchor)
				t.RequestRepaint()
				return true
			}
		}

	case event.EventMouseMove:
		if t.isDragging {
			if data, ok := evt.Data.(event.MouseEvent); ok {
				t.pendingMouseX = data.X
				t.RequestRepaint()
				return true
			}
		}

	case event.EventMouseRelease:
		if t.isDragging {
			t.isDragging = false
			t.RequestRepaint()
			return true
		}

	case event.EventKeyPress:
		if t.isFocused {
			if data, ok := evt.Data.(event.KeyEvent); ok {
				runes := []rune(t.Text)
				isShift := (data.Modifiers & event.ModShift) != 0

				// Navigation keys are allowed in ReadOnly mode
				switch data.VirtualKeyCode {
				case 0x25: // VK_LEFT
					if t.cursorPos > 0 {
						t.cursorPos--
					}
					if !isShift {
						t.selStart = t.cursorPos
					}
					t.RequestRepaint()
					return true
				case 0x27: // VK_RIGHT
					if t.cursorPos < len(runes) {
						t.cursorPos++
					}
					if !isShift {
						t.selStart = t.cursorPos
					}
					t.RequestRepaint()
					return true
				}

				if t.ReadOnly {
					return false
				}

				switch data.VirtualKeyCode {
				case 0x08: // VK_BACK
					if len(runes) > 0 {
						if t.selStart != t.cursorPos {
							// Delete selection
							start, end := t.selStart, t.cursorPos
							if start > end {
								start, end = end, start
							}
							t.Text = string(runes[:start]) + string(runes[end:])
							t.cursorPos = start
							t.selStart = start
						} else if t.cursorPos > 0 {
							t.Text = string(runes[:t.cursorPos-1]) + string(runes[t.cursorPos:])
							t.cursorPos--
							t.selStart = t.cursorPos
						}
					}
					t.RequestRepaint()
					return true
				}
			}
		}

	case event.EventChar:
		if t.isFocused && !t.ReadOnly {
			if data, ok := evt.Data.(event.KeyEvent); ok {
				// Ignore control chars
				if data.Rune < 32 {
					return false
				}

				runes := []rune(t.Text)

				// If selection exists, replace it
				if t.selStart != t.cursorPos {
					start, end := t.selStart, t.cursorPos
					if start > end {
						start, end = end, start
					}
					// Remove selection
					t.Text = string(runes[:start]) + string(data.Rune) + string(runes[end:])
					t.cursorPos = start + 1
					t.selStart = t.cursorPos
				} else {
					// Insert at cursor
					left := runes[:t.cursorPos]
					right := runes[t.cursorPos:]
					t.Text = string(left) + string(data.Rune) + string(right)
					t.cursorPos++
					t.selStart = t.cursorPos
				}

				t.RequestRepaint()
				return true
			}
		}
	}
	return false
}
