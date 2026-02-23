package main

import (
	"fmt"
	"github.com/jacksalad/goui_v0/component"
	"github.com/jacksalad/goui_v0/layout"
	"github.com/jacksalad/goui_v0/window"
	"runtime"
)

func main() {
	runtime.LockOSThread()

	win, err := window.NewWindow(window.WindowConfig{
		Title:     "Layout Example - Flex & Grid",
		Width:     800,
		Height:    600,
		Resizable: true,
	})
	if err != nil {
		panic(err)
	}

	// Root uses VBox to stack Flex and Grid demos
	win.Root.SetLayout(&layout.VBoxLayout{Spacing: 20, Padding: 10})

	// 1. Flex Layout Demo (Row)
	flexLabel := component.NewLabel("Flex Layout (Row): Fixed, Grow 1, Grow 2")
	win.Root.Add(flexLabel)

	flexPanel := component.NewPanel(0, 0, 0, 80) // Height 80
	flexPanel.BgColor = 0xFFE0E0E0
	
	// Create FlexLayout
	flex := layout.NewFlexLayout(layout.FlexRow)
	flex.JustifyContent = layout.FlexStart
	flex.AlignItems = layout.AlignCenter
	flex.Padding = 10
	flex.Spacing = 10
	flexPanel.SetLayout(flex)
	
	// Add buttons
	btn1 := component.NewButton("Fixed")
	flexPanel.Add(btn1)
	
	btn2 := component.NewButton("Grow 1")
	flexPanel.Add(btn2)
	flex.SetGrow(btn2, 1)
	
	btn3 := component.NewButton("Grow 2")
	flexPanel.Add(btn3)
	flex.SetGrow(btn3, 2)
	
	win.Root.Add(flexPanel)

	// 2. Grid Layout Demo
	gridLabel := component.NewLabel("Grid Layout (3x3):")
	win.Root.Add(gridLabel)

	gridPanel := component.NewPanel(0, 0, 0, 300) // Height 300
	gridPanel.BgColor = 0xFFD0D0D0
	
	// Create GridLayout 3x3
	grid := layout.NewGridLayout(3, 3)
	grid.Padding = 10
	grid.Spacing = 5
	gridPanel.SetLayout(grid)
	
	for i := 0; i < 9; i++ {
		l := component.NewButton(fmt.Sprintf("Cell %d", i))
		gridPanel.Add(l)
		row := i / 3
		col := i % 3
		grid.SetPosition(l, row, col)
	}
	
	// Add a spanning item (overwrites some cells visually)
	spanBtn := component.NewButton("Span 2 Cols")
	// Let's put it at row 1, col 1, spanning 2 cols
	gridPanel.Add(spanBtn)
	grid.SetPosition(spanBtn, 1, 1)
	grid.SetSpan(spanBtn, 1, 2)
	
	win.Root.Add(gridPanel)

	win.Show()
	win.Run()
}
