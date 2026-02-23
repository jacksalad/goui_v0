package main

import (
	"goui/component"
	"goui/layout"
	"goui/render"
	"goui/window"
	"strconv"
)

type Calculator struct {
	window  *window.Window
	display *component.TextBox

	// State
	currentVal  float64
	storedVal   float64
	operation   string
	shouldReset bool
}

func NewCalculator() (*Calculator, error) {
	win, err := window.NewWindow(window.WindowConfig{
		Title:     "GOUI Calculator",
		Width:     320,
		Height:    470,
		Resizable: false,
	})
	if err != nil {
		return nil, err
	}

	calc := &Calculator{
		window:      win,
		shouldReset: false,
	}

	calc.setupUI()
	return calc, nil
}

func (c *Calculator) setupUI() {
	// Main Layout: VBox
	rootLayout := &layout.VBoxLayout{
		Padding: 10,
		Spacing: 10,
	}
	c.window.Root.SetLayout(rootLayout)
	c.window.Root.BgColor = 0xFFF0F0F0

	// Fonts
	displayFont := render.NewFont("SimHei", 24)
	btnFont := render.NewFont("SimHei", 16)

	// Display Area
	c.display = component.NewTextBox(300)
	c.display.Text = "0"
	c.display.ReadOnly = true
	c.display.Font = displayFont
	c.window.Root.Add(c.display)

	// Button Grid
	gridPanel := component.NewPanel(0, 0, 300, 350)
	gridLayout := layout.NewGridLayout(5, 4) // 5 rows, 4 cols
	gridLayout.Spacing = 5
	gridPanel.SetLayout(gridLayout)

	c.window.Root.Add(gridPanel)

	// Buttons definition
	// Row 0: C, /, *, -
	// Row 1: 7, 8, 9, +
	// Row 2: 4, 5, 6, (Span + down?) No, let's keep it simple
	// Row 3: 1, 2, 3, = (Span 2 rows?)
	// Row 4: 0 (Span 2), .

	// Let's use a standard grid layout
	buttons := []struct {
		text             string
		row, col         int
		rowSpan, colSpan int
		action           func()
	}{
		// Row 0
		{"C", 0, 0, 1, 1, c.onClear},
		{"/", 0, 1, 1, 1, func() { c.onOp("/") }},
		{"*", 0, 2, 1, 1, func() { c.onOp("*") }},
		{"-", 0, 3, 1, 1, func() { c.onOp("-") }},

		// Row 1
		{"7", 1, 0, 1, 1, func() { c.onDigit("7") }},
		{"8", 1, 1, 1, 1, func() { c.onDigit("8") }},
		{"9", 1, 2, 1, 1, func() { c.onDigit("9") }},
		{"+", 1, 3, 2, 1, func() { c.onOp("+") }}, // Span 2 rows?

		// Row 2
		{"4", 2, 0, 1, 1, func() { c.onDigit("4") }},
		{"5", 2, 1, 1, 1, func() { c.onDigit("5") }},
		{"6", 2, 2, 1, 1, func() { c.onDigit("6") }},
		// + spans here

		// Row 3
		{"1", 3, 0, 1, 1, func() { c.onDigit("1") }},
		{"2", 3, 1, 1, 1, func() { c.onDigit("2") }},
		{"3", 3, 2, 1, 1, func() { c.onDigit("3") }},
		{"=", 3, 3, 2, 1, c.onEqual}, // Span 2 rows

		// Row 4
		{"0", 4, 0, 1, 2, func() { c.onDigit("0") }}, // Span 2 cols
		// 0 spans here
		{".", 4, 2, 1, 1, c.onDot},
		// = spans here
	}

	for _, b := range buttons {
		btn := component.NewButton(b.text)
		btn.OnClick = b.action
		btn.Font = btnFont
		gridPanel.Add(btn)
		gridLayout.SetPosition(btn, b.row, b.col)
		if b.rowSpan > 1 || b.colSpan > 1 {
			gridLayout.SetSpan(btn, b.rowSpan, b.colSpan)
		}
	}
}

func (c *Calculator) updateDisplay(val string) {
	c.display.Text = val
	c.display.RequestRepaint()
}

func (c *Calculator) onDigit(digit string) {
	if c.shouldReset {
		c.display.Text = ""
		c.shouldReset = false
	}

	if c.display.Text == "0" {
		c.display.Text = digit
	} else {
		c.display.Text += digit
	}
}

func (c *Calculator) onDot() {
	if c.shouldReset {
		c.display.Text = "0."
		c.shouldReset = false
		return
	}

	// Check if dot already exists
	for _, r := range c.display.Text {
		if r == '.' {
			return
		}
	}
	c.display.Text += "."
}

func (c *Calculator) onClear() {
	c.display.Text = "0"
	c.storedVal = 0
	c.operation = ""
	c.shouldReset = false
}

func (c *Calculator) onOp(op string) {
	// If we have an operation pending, calculate it first?
	// For simple calc, just store current value
	val, _ := strconv.ParseFloat(c.display.Text, 64)
	c.storedVal = val
	c.operation = op
	c.shouldReset = true
}

func (c *Calculator) onEqual() {
	if c.operation == "" {
		return
	}

	val, _ := strconv.ParseFloat(c.display.Text, 64)
	var result float64

	switch c.operation {
	case "+":
		result = c.storedVal + val
	case "-":
		result = c.storedVal - val
	case "*":
		result = c.storedVal * val
	case "/":
		if val != 0 {
			result = c.storedVal / val
		} else {
			c.display.Text = "Error"
			c.shouldReset = true
			return
		}
	}

	// Format result
	c.display.Text = strconv.FormatFloat(result, 'f', -1, 64)
	c.shouldReset = true
	c.operation = ""
}

func (c *Calculator) Run() {
	c.window.Show()
	c.window.Run()
}

func main() {
	calc, err := NewCalculator()
	if err != nil {
		panic(err)
	}
	calc.Run()
}
