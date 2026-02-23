package main

import (
	"fmt"
	"goui/component"
	"goui/layout"
	"goui/render"
	"goui/window"
)

type TicTacToe struct {
	window    *window.Window
	status    *component.Label
	buttons   [9]*component.Button
	board     [9]string // "" for empty, "X", "O"
	turn      string    // "X" or "O"
	gameOver  bool
	moveCount int

	// Fonts
	titleFont  *render.Font
	gameFont   *render.Font
	buttonFont *render.Font
}

func NewTicTacToe() (*TicTacToe, error) {
	win, err := window.NewWindow(window.WindowConfig{
		Title:     "Tic-Tac-Toe",
		Width:     320,
		Height:    450,
		Resizable: false,
	})
	if err != nil {
		return nil, err
	}

	game := &TicTacToe{
		window: win,
		turn:   "X",
	}
	game.setupUI()
	return game, nil
}

func (g *TicTacToe) setupUI() {
	// Fonts
	g.titleFont = render.NewFont("SimHei", 20)
	g.gameFont = render.NewFont("SimHei", 16)
	g.buttonFont = render.NewFont("Arial", 32) // Larger for X/O

	// Main Layout
	rootLayout := &layout.VBoxLayout{
		Padding: 10,
		Spacing: 10,
	}
	g.window.Root.SetLayout(rootLayout)
	g.window.Root.BgColor = 0xFFF0F0F0

	// Status Label
	g.status = component.NewLabel("Player X's Turn")
	g.status.Font = g.titleFont
	g.window.Root.Add(g.status)

	// Game Grid Panel
	// Height 300 is enough for 3 rows of buttons
	gridPanel := component.NewPanel(0, 0, 300, 300)
	gridLayout := layout.NewGridLayout(3, 3)
	gridLayout.Spacing = 5
	gridLayout.Padding = 0

	// We need to set the layout to the panel
	// Note: We'll set it after adding children so they are arranged,
	// but SetLayout calls LayoutChildren automatically.
	// However, we need the GridLayout instance to set positions.
	gridPanel.SetLayout(gridLayout)

	g.window.Root.Add(gridPanel)

	// Buttons
	for i := 0; i < 9; i++ {
		idx := i
		btn := component.NewButton("")
		btn.Font = g.buttonFont
		btn.OnClick = func() {
			g.onCellClick(idx)
		}
		g.buttons[i] = btn
		gridPanel.Add(btn)

		row := i / 3
		col := i % 3
		gridLayout.SetPosition(btn, row, col)
	}

	// Reset Button
	resetBtn := component.NewButton("Reset Game")
	resetBtn.Font = g.gameFont
	resetBtn.OnClick = g.resetGame
	g.window.Root.Add(resetBtn)
}

func (g *TicTacToe) onCellClick(idx int) {
	if g.gameOver || g.board[idx] != "" {
		return
	}

	// Update Board
	g.board[idx] = g.turn
	g.buttons[idx].Text = g.turn
	g.buttons[idx].RequestRepaint()
	g.moveCount++

	// Check Win
	if g.checkWin() {
		g.status.Text = fmt.Sprintf("Player %s Wins!", g.turn)
		g.gameOver = true
	} else if g.moveCount == 9 {
		g.status.Text = "Draw!"
		g.gameOver = true
	} else {
		// Switch Turn
		if g.turn == "X" {
			g.turn = "O"
		} else {
			g.turn = "X"
		}
		g.status.Text = fmt.Sprintf("Player %s's Turn", g.turn)
	}
	g.status.RequestRepaint()
}

func (g *TicTacToe) checkWin() bool {
	// Rows
	for i := 0; i < 3; i++ {
		if g.checkLine(i*3, i*3+1, i*3+2) {
			return true
		}
	}
	// Cols
	for i := 0; i < 3; i++ {
		if g.checkLine(i, i+3, i+6) {
			return true
		}
	}
	// Diagonals
	if g.checkLine(0, 4, 8) {
		return true
	}
	if g.checkLine(2, 4, 6) {
		return true
	}
	return false
}

func (g *TicTacToe) checkLine(a, b, c int) bool {
	return g.board[a] != "" && g.board[a] == g.board[b] && g.board[b] == g.board[c]
}

func (g *TicTacToe) resetGame() {
	g.turn = "X"
	g.gameOver = false
	g.moveCount = 0
	g.status.Text = "Player X's Turn"
	g.status.RequestRepaint()

	for i := 0; i < 9; i++ {
		g.board[i] = ""
		g.buttons[i].Text = ""
		g.buttons[i].RequestRepaint()
	}
}

func main() {
	game, err := NewTicTacToe()
	if err != nil {
		panic(err)
	}
	game.window.Show()
	game.window.Run()
}
