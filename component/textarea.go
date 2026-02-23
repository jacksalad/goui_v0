package component

import (
	"github.com/jacksalad/goui_v0/event"
	"github.com/jacksalad/goui_v0/render"
	"strings"
	"time"
)

type TextArea struct {
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
	pendingMouseX int32
	pendingMouseY int32
	cursorBlink   bool
	lastBlink     int64

	// Scroll (simple)
	scrollY    int32
	lineHeight int32
}

func NewTextArea(width, height int32) *TextArea {
	t := &TextArea{
		cursorPos:  0,
		selStart:   -1, // -1 means no selection anchor
		lineHeight: 20, // Default estimate
	}
	t.SetBounds(0, 0, width, height)
	t.Visible = true
	return t
}

func (t *TextArea) GetPreferredSize() (int32, int32) {
	// For TextArea, preferred size is usually its bounds or content size
	// Let's return current bounds for now
	return t.Bounds.Width, t.Bounds.Height
}

func (t *TextArea) getLineHeight(canvas *render.Canvas) int32 {
	if t.Font != nil {
		canvas.SetFont(t.Font)
	}
	_, h := canvas.MeasureText("Tg")
	t.lineHeight = h + 4
	return t.lineHeight
}

func (t *TextArea) ensureCursorVisible() {
	line, _ := t.getLineCol(t.cursorPos)

	// Padding
	paddingY := int32(5)

	// Calculate target Y of cursor top and bottom
	cursorTop := paddingY + int32(line)*t.lineHeight
	cursorBottom := cursorTop + t.lineHeight

	// Current View
	viewTop := t.scrollY
	viewBottom := t.scrollY + t.Bounds.Height

	// Scroll if needed
	if cursorTop < viewTop {
		t.scrollY = cursorTop
	} else if cursorBottom > viewBottom {
		t.scrollY = cursorBottom - t.Bounds.Height
	}

	// Clamp scrollY
	if t.scrollY < 0 {
		t.scrollY = 0
	}
}

// Helper to convert linear position to (line, col)
func (t *TextArea) getLineCol(pos int) (int, int) {
	runes := []rune(t.Text)
	if pos > len(runes) {
		pos = len(runes)
	}

	line := 0
	col := 0
	for i := 0; i < pos; i++ {
		if runes[i] == '\n' {
			line++
			col = 0
		} else {
			col++
		}
	}
	return line, col
}

// Helper to convert (line, col) to linear position
func (t *TextArea) getPosFromLineCol(line, col int) int {
	lines := strings.Split(t.Text, "\n")
	if line < 0 {
		return 0
	}
	if line >= len(lines) {
		return len([]rune(t.Text))
	}

	pos := 0
	for i := 0; i < line; i++ {
		pos += len([]rune(lines[i])) + 1 // +1 for \n
	}

	lineLen := len([]rune(lines[line]))
	if col > lineLen {
		col = lineLen
	}
	return pos + col
}

func (t *TextArea) getIndexFromXY(x, y int32, canvas *render.Canvas) int {
	// Adjust for padding/scroll
	localX := x - (t.Bounds.X + 5)
	localY := y - (t.Bounds.Y + 5) + t.scrollY

	if localY < 0 {
		return 0
	}

	lineHeight := t.getLineHeight(canvas)
	lineIndex := int(localY / lineHeight)

	lines := strings.Split(t.Text, "\n")
	if lineIndex >= len(lines) {
		return len([]rune(t.Text))
	}

	// Find col in line
	lineStr := lines[lineIndex]
	runes := []rune(lineStr)

	// Linear search for X
	bestCol := 0
	bestDist := int32(999999)

	for i := 0; i <= len(runes); i++ {
		w, _ := canvas.MeasureText(string(runes[:i]))
		dist := abs(w - localX)
		if dist < bestDist {
			bestDist = dist
			bestCol = i
		}
	}

	// Calculate global pos
	pos := 0
	for i := 0; i < lineIndex; i++ {
		pos += len([]rune(lines[i])) + 1 // +1 for \n
	}
	return pos + bestCol
}

func abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func (t *TextArea) Render(canvas *render.Canvas) {
	if !t.Visible {
		return
	}

	if t.Font != nil {
		canvas.SetFont(t.Font)
	}

	// Check for pending mouse click (needs canvas for measurement)
	if t.pendingMouseX != 0 || t.pendingMouseY != 0 {
		t.cursorPos = t.getIndexFromXY(t.pendingMouseX, t.pendingMouseY, canvas)
		t.pendingMouseX = 0
		t.pendingMouseY = 0
	}

	// Background
	canvas.FillRect(t.Bounds.X, t.Bounds.Y, t.Bounds.Width, t.Bounds.Height, 0xFFFFFFFF)
	// Border
	borderColor := uint32(0xFFAAAAAA)
	if t.isFocused {
		borderColor = 0xFF0078D7
	}
	canvas.FillRect(t.Bounds.X, t.Bounds.Y, t.Bounds.Width, 1, borderColor)
	canvas.FillRect(t.Bounds.X, t.Bounds.Y+t.Bounds.Height-1, t.Bounds.Width, 1, borderColor)
	canvas.FillRect(t.Bounds.X, t.Bounds.Y, 1, t.Bounds.Height, borderColor)
	canvas.FillRect(t.Bounds.X+t.Bounds.Width-1, t.Bounds.Y, 1, t.Bounds.Height, borderColor)

	// Draw Text
	lines := strings.Split(t.Text, "\n")
	lineHeight := t.getLineHeight(canvas)

	paddingX := int32(5)
	paddingY := int32(5)

	startX := t.Bounds.X + paddingX
	startY := t.Bounds.Y + paddingY - t.scrollY

	for i, line := range lines {
		y := startY + int32(i)*lineHeight
		if y+lineHeight < t.Bounds.Y {
			continue // above view
		}
		if y > t.Bounds.Y+t.Bounds.Height {
			break // below view
		}

		canvas.DrawText(startX, y, line, 0xFF000000)
	}

	// Draw Cursor
	if t.isFocused {
		now := time.Now().UnixNano()
		if now-t.lastBlink > 500*1e6 {
			t.cursorBlink = !t.cursorBlink
			t.lastBlink = now
			t.RequestRepaint()
		}

		if t.cursorBlink {
			line, col := t.getLineCol(t.cursorPos)
			if line < len(lines) {
				lineStr := lines[line]
				runes := []rune(lineStr)
				if col > len(runes) {
					col = len(runes)
				}
				w, _ := canvas.MeasureText(string(runes[:col]))

				cursorX := startX + w
				cursorY := startY + int32(line)*lineHeight

				// Ensure cursor inside bounds
				if cursorY >= t.Bounds.Y && cursorY+lineHeight <= t.Bounds.Y+t.Bounds.Height {
					canvas.FillRect(cursorX, cursorY, 2, lineHeight, 0xFF000000)
				}
			}
		}
	}
	t.RepaintRequested = false
}

func (t *TextArea) OnFocus() {
	t.isFocused = true
	t.RequestRepaint()
}

func (t *TextArea) OnBlur() {
	t.isFocused = false
	t.RequestRepaint()
}

func (t *TextArea) OnEvent(evt event.Event) bool {
	if !t.Visible {
		return false
	}

	switch evt.Type {
	case event.EventMouseClick:
		if data, ok := evt.Data.(event.MouseEvent); ok {
			if t.Bounds.Contains(data.X, data.Y) {
				t.isFocused = true
				t.pendingMouseX = data.X
				t.pendingMouseY = data.Y
				t.RequestRepaint()
				return true
			} else {
				if t.isFocused {
					t.isFocused = false
					t.RequestRepaint()
				}
			}
		}

	case event.EventMouseWheel:
		if data, ok := evt.Data.(event.MouseEvent); ok {
			// Scroll
			scrollAmount := int32(data.Delta) * -1 // Windows delta is + for up, usually we scroll up (decrease Y)
			// But wait, scrollY is offset. Increasing scrollY moves content UP.
			// So if wheel is UP (+120), we want to DECREASE scrollY (show content above).
			// So scrollAmount should be negative.

			// Adjust speed
			step := int32(40)
			if scrollAmount > 0 {
				t.scrollY += step
			} else {
				t.scrollY -= step
			}

			if t.scrollY < 0 {
				t.scrollY = 0
			}
			// Max scroll?
			// Calculate total height
			lines := strings.Split(t.Text, "\n")
			totalHeight := int32(len(lines))*t.lineHeight + 10 // padding
			if t.scrollY > totalHeight-t.Bounds.Height {
				t.scrollY = totalHeight - t.Bounds.Height
			}
			if t.scrollY < 0 {
				t.scrollY = 0
			}

			t.RequestRepaint()
			return true
		}

	case event.EventKeyPress:
		if t.isFocused {
			if data, ok := evt.Data.(event.KeyEvent); ok {
				runes := []rune(t.Text)
				if t.cursorPos > len(runes) {
					t.cursorPos = len(runes)
				}
				if t.cursorPos < 0 {
					t.cursorPos = 0
				}

				switch data.VirtualKeyCode {
				case 8: // Backspace
					if t.cursorPos > 0 {
						t.Text = string(runes[:t.cursorPos-1]) + string(runes[t.cursorPos:])
						t.cursorPos--
						t.ensureCursorVisible()
						t.RequestRepaint()
					}
					return true

				case 37: // Left
					if t.cursorPos > 0 {
						t.cursorPos--
						t.ensureCursorVisible()
						t.RequestRepaint()
					}
					return true

				case 39: // Right
					if t.cursorPos < len(runes) {
						t.cursorPos++
						t.ensureCursorVisible()
						t.RequestRepaint()
					}
					return true

				case 38: // Up
					line, col := t.getLineCol(t.cursorPos)
					if line > 0 {
						t.cursorPos = t.getPosFromLineCol(line-1, col)
						t.ensureCursorVisible()
						t.RequestRepaint()
					}
					return true

				case 40: // Down
					line, col := t.getLineCol(t.cursorPos)
					lines := strings.Split(t.Text, "\n")
					if line < len(lines)-1 {
						t.cursorPos = t.getPosFromLineCol(line+1, col)
						t.ensureCursorVisible()
						t.RequestRepaint()
					}
					return true

				case 13: // Enter
					t.Text = string(runes[:t.cursorPos]) + "\n" + string(runes[t.cursorPos:])
					t.cursorPos++
					t.ensureCursorVisible()
					t.RequestRepaint()
					return true
				}
			}
		}

	case event.EventChar:
		if t.isFocused {
			if data, ok := evt.Data.(event.KeyEvent); ok {
				runes := []rune(t.Text)
				if t.cursorPos > len(runes) {
					t.cursorPos = len(runes)
				}
				if t.cursorPos < 0 {
					t.cursorPos = 0
				}

				// Filter control chars (BS=8, CR=13, etc.)
				if data.Rune >= 32 {
					runes := []rune(t.Text)
					t.Text = string(runes[:t.cursorPos]) + string(data.Rune) + string(runes[t.cursorPos:])
					t.cursorPos++
					t.ensureCursorVisible()
					t.RequestRepaint()
					return true
				}
			}
		}
	}
	return false
}
