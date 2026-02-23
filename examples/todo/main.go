package main

import (
	"fmt"
	"github.com/jacksalad/goui_v0/component"
	"github.com/jacksalad/goui_v0/layout"
	"github.com/jacksalad/goui_v0/render"
	"github.com/jacksalad/goui_v0/window"
)

type TodoApp struct {
	window    *window.Window
	input     *component.TextBox
	listPanel *component.Panel
	
	// Fonts
	titleFont *render.Font
	itemFont  *render.Font
}

func NewTodoApp() (*TodoApp, error) {
	win, err := window.NewWindow(window.WindowConfig{
		Title:     "GOUI Todo List",
		Width:     400,
		Height:    600,
		Resizable: false,
	})
	if err != nil {
		return nil, err
	}

	app := &TodoApp{
		window: win,
	}
	app.setupUI()
	return app, nil
}

func (app *TodoApp) setupUI() {
	// Fonts
	app.titleFont = render.NewFont("SimHei", 24)
	app.itemFont = render.NewFont("SimHei", 16)

	// Main Layout (VBox)
	rootLayout := &layout.VBoxLayout{
		Padding: 20,
		Spacing: 10,
	}
	app.window.Root.SetLayout(rootLayout)
	app.window.Root.BgColor = 0xFFF5F5F5

	// Title
	title := component.NewLabel("Todo List")
	title.FgColor = 0xFF333333
	title.Font = app.titleFont
	app.window.Root.Add(title)

	// Input Area (Panel with HBox)
	inputPanel := component.NewPanel(0, 0, 360, 60)
	inputLayout := &layout.HBoxLayout{
		Padding: 0,
		Spacing: 10,
	}
	inputPanel.SetLayout(inputLayout)
	inputPanel.BgColor = 0xFFF5F5F5 // Match background

	// Input Box
	app.input = component.NewTextBox(260)
	app.input.Font = app.itemFont
	app.input.Placeholder = "Enter new task..."
	inputPanel.Add(app.input)

	// Add Button
	addBtn := component.NewButton("Add")
	addBtn.Font = app.itemFont
	addBtn.OnClick = app.addTodo
	inputPanel.Add(addBtn)

	app.window.Root.Add(inputPanel)

	// Todo List Container
	app.listPanel = component.NewPanel(0, 0, 360, 400)
	listLayout := &layout.VBoxLayout{
		Padding: 10,
		Spacing: 5,
	}
	app.listPanel.SetLayout(listLayout)
	app.listPanel.BgColor = 0xFFFFFFFF // White background for list
	
	app.window.Root.Add(app.listPanel)
}

func (app *TodoApp) addTodo() {
	text := app.input.Text
	if text == "" {
		return
	}

	// Create Todo Item Panel (HBox)
	// Using a panel to group checkbox and delete button
	itemPanel := component.NewPanel(0, 0, 340, 40)
	itemLayout := &layout.HBoxLayout{
		Padding: 5,
		Spacing: 10,
	}
	itemPanel.SetLayout(itemLayout)
	itemPanel.BgColor = 0xFFF0F8FF // Light Alice Blue

	// Checkbox
	chk := component.NewCheckBox(text)
	chk.Font = app.itemFont
	chk.OnCheck = func(checked bool) {
		fmt.Printf("Task '%s' checked: %v\n", text, checked)
	}
	
	// Delete Button
	delBtn := component.NewButton("X")
	delBtn.Font = app.itemFont
	// Small red button? Button doesn't support color yet, but we can imagine
	delBtn.OnClick = func() {
		app.listPanel.Remove(itemPanel)
		// Force repaint
		app.window.Root.RequestRepaint()
		// Ideally layout should update automatically
	}

	itemPanel.Add(chk)
	itemPanel.Add(delBtn)

	app.listPanel.Add(itemPanel)
	
	// Clear input
	app.input.Text = ""
	app.input.RequestRepaint()
	
	// Force repaint of list
	app.listPanel.RequestRepaint()
	app.window.Root.RequestRepaint()
}

func main() {
	app, err := NewTodoApp()
	if err != nil {
		panic(err)
	}
	app.window.Show()
	app.window.Run()
}
