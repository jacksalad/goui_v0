package main

import (
	"fmt"
	"goui/component"
	"goui/event"
	"goui/layout"
	"goui/render"
	"goui/window"
	"math/rand"
	"sync"
	"time"
)

// Directions (using int for simplicity)
const (
	DirUp = iota
	DirDown
	DirLeft
	DirRight
)

const (
	GridSize   = 20
	GridWidth  = 30
	GridHeight = 20
)

type Point struct {
	X, Y int
}

type SnakeGame struct {
	component.BaseComponent

	snake     []Point
	food      Point
	direction int
	nextDir   int // Buffer next direction to prevent 180 turns in one frame
	gameOver  bool
	score     int
	paused    bool

	mu sync.Mutex // Protect game state from concurrent access (Render vs GameLoop)
}

func NewSnakeGame() *SnakeGame {
	g := &SnakeGame{
		snake:     []Point{{5, 5}, {4, 5}, {3, 5}},
		food:      Point{15, 10},
		direction: DirRight,
		nextDir:   DirRight,
	}
	g.Visible = true
	return g
}

func (g *SnakeGame) GetPreferredSize() (int32, int32) {
	return int32(GridWidth * GridSize), int32(GridHeight * GridSize)
}

func (g *SnakeGame) Render(canvas *render.Canvas) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Draw Background
	canvas.FillRect(g.Bounds.X, g.Bounds.Y, g.Bounds.Width, g.Bounds.Height, 0xFF222222) // Dark Gray

	if g.gameOver {
		// Draw Game Over Text
		text := fmt.Sprintf("GAME OVER! Score: %d", g.score)
		// Simple centering logic (approximate)
		canvas.DrawText(g.Bounds.X+50, g.Bounds.Y+g.Bounds.Height/2-10, text, 0xFFFFFFFF)
		canvas.DrawText(g.Bounds.X+50, g.Bounds.Y+g.Bounds.Height/2+20, "Press ENTER to Restart", 0xFFAAAAAA)
		return
	}

	// Draw Food
	canvas.FillRect(
		g.Bounds.X+int32(g.food.X*GridSize),
		g.Bounds.Y+int32(g.food.Y*GridSize),
		GridSize-1,
		GridSize-1,
		0xFFFF0000, // Red
	)

	// Draw Snake
	for i, p := range g.snake {
		color := uint32(0xFF00FF00) // Green
		if i == 0 {
			color = 0xFF00CC00 // Darker Green for Head
		}

		canvas.FillRect(
			g.Bounds.X+int32(p.X*GridSize),
			g.Bounds.Y+int32(p.Y*GridSize),
			GridSize-1,
			GridSize-1,
			color,
		)
	}

	// Draw Pause Overlay
	if g.paused {
		canvas.DrawText(g.Bounds.X+10, g.Bounds.Y+10, "PAUSED", 0xFFFFFF00)
	}
}

func (g *SnakeGame) OnEvent(evt event.Event) bool {
	if evt.Type == event.EventKeyPress {
		key := evt.Data.(event.KeyEvent)

		g.mu.Lock()
		defer g.mu.Unlock()

		if g.gameOver {
			if key.VirtualKeyCode == 0x0D { // Enter
				g.restart()
				return true
			}
			return false
		}

		switch key.VirtualKeyCode {
		case 0x26: // Up Arrow
			if g.direction != DirDown {
				g.nextDir = DirUp
			}
		case 0x28: // Down Arrow
			if g.direction != DirUp {
				g.nextDir = DirDown
			}
		case 0x25: // Left Arrow
			if g.direction != DirRight {
				g.nextDir = DirLeft
			}
		case 0x27: // Right Arrow
			if g.direction != DirLeft {
				g.nextDir = DirRight
			}
		case 0x20: // Space
			g.paused = !g.paused
		}
		return true
	}
	return false
}

func (g *SnakeGame) restart() {
	g.snake = []Point{{5, 5}, {4, 5}, {3, 5}}
	g.food = Point{15, 10}
	g.direction = DirRight
	g.nextDir = DirRight
	g.gameOver = false
	g.score = 0
	g.paused = false
}

func (g *SnakeGame) Step() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.gameOver || g.paused {
		return
	}

	g.direction = g.nextDir

	// Calculate new head position
	head := g.snake[0]
	newHead := head

	switch g.direction {
	case DirUp:
		newHead.Y--
	case DirDown:
		newHead.Y++
	case DirLeft:
		newHead.X--
	case DirRight:
		newHead.X++
	}

	// Check Collisions (Walls)
	if newHead.X < 0 || newHead.X >= GridWidth || newHead.Y < 0 || newHead.Y >= GridHeight {
		g.gameOver = true
		return
	}

	// Check Collisions (Self)
	for _, p := range g.snake {
		if p == newHead {
			g.gameOver = true
			return
		}
	}

	// Move Snake
	g.snake = append([]Point{newHead}, g.snake...)

	// Check Food
	if newHead == g.food {
		g.score += 10
		// Spawn new food
		for {
			g.food = Point{rand.Intn(GridWidth), rand.Intn(GridHeight)}
			// Make sure food doesn't spawn on snake
			onSnake := false
			for _, p := range g.snake {
				if p == g.food {
					onSnake = true
					break
				}
			}
			if !onSnake {
				break
			}
		}
	} else {
		// Remove tail
		g.snake = g.snake[:len(g.snake)-1]
	}
}

func main() {
	// 1. Create Window
	win, err := window.NewWindow(window.WindowConfig{
		Title:     "GoUI Snake",
		Width:     800,
		Height:    600,
		Resizable: false,
	})
	if err != nil {
		panic(err)
	}

	// 2. Setup Layout
	win.Root.BgColor = 0xFF333333
	win.Root.SetLayout(&layout.VBoxLayout{
		Padding: 20,
		Spacing: 20,
	})

	// 3. Add UI Elements
	// Title / Score Panel
	header := component.NewPanel(0, 0, 0, 50)
	header.BgColor = 0xFF444444
	header.SetLayout(&layout.HBoxLayout{
		Padding: 10,
		Spacing: 20,
	})
	win.Add(header)

	title := component.NewLabel("GoUI Snake Game")
	title.FgColor = 0xFFFFFFFF
	header.Add(title)

	scoreLbl := component.NewLabel("Score: 0")
	scoreLbl.FgColor = 0xFF00FF00
	header.Add(scoreLbl)

	helpLbl := component.NewLabel("Use Arrow Keys to Move | Space to Pause")
	helpLbl.FgColor = 0xFFAAAAAA
	header.Add(helpLbl)

	// Game Area
	game := NewSnakeGame()
	// Center the game board roughly by putting it in a container or just add it
	// Since we use VBox, it will be added below the header
	win.Add(game)

	// Focus the game component to receive keyboard events
	win.SetFocus(game)

	// 4. Game Loop
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond) // 10 FPS
		defer ticker.Stop()

		for range ticker.C {
			game.Step()

			// Update Score Label
			// Note: Accessing UI components from another goroutine is generally unsafe
			// if they are not protected. Here we just set a string, which might race
			// if the renderer reads it at the exact same time.
			// Ideally we should use a safe way, but for this demo:
			scoreLbl.Text = fmt.Sprintf("Score: %d", game.score)

			win.RequestRepaint()
		}
	}()

	win.Show()
	win.Run()
}
