package window

import (
	"goui/component"
	"goui/event"
	"goui/render"
	"runtime"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	moduser32   = windows.NewLazySystemDLL("user32.dll")
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procGetModuleHandleW = modkernel32.NewProc("GetModuleHandleW")
	procRegisterClassExW = moduser32.NewProc("RegisterClassExW")
	procCreateWindowExW  = moduser32.NewProc("CreateWindowExW")
	procDefWindowProcW   = moduser32.NewProc("DefWindowProcW")
	procShowWindow       = moduser32.NewProc("ShowWindow")
	procUpdateWindow     = moduser32.NewProc("UpdateWindow")
	procGetMessageW      = moduser32.NewProc("GetMessageW")
	procTranslateMessage = moduser32.NewProc("TranslateMessage")
	procDispatchMessageW = moduser32.NewProc("DispatchMessageW")
	procPostQuitMessage  = moduser32.NewProc("PostQuitMessage")
	procPostMessageW     = moduser32.NewProc("PostMessageW")
	procLoadCursorW      = moduser32.NewProc("LoadCursorW")
	procDestroyWindow    = moduser32.NewProc("DestroyWindow")
	procSetTimer         = moduser32.NewProc("SetTimer")
	procGetKeyState      = moduser32.NewProc("GetKeyState")
)

var (
	windowsMap = make(map[windows.Handle]*Window)
	mapMu      sync.RWMutex
)

const (
	IDC_ARROW           = 32512
	CS_HREDRAW          = 0x0002
	CS_VREDRAW          = 0x0001
	COLOR_WINDOW        = 5
	WS_OVERLAPPED       = 0x00000000
	WS_CAPTION          = 0x00C00000
	WS_SYSMENU          = 0x00080000
	WS_THICKFRAME       = 0x00040000
	WS_MINIMIZEBOX      = 0x00020000
	WS_MAXIMIZEBOX      = 0x00010000
	WS_OVERLAPPEDWINDOW = WS_OVERLAPPED | WS_CAPTION | WS_SYSMENU | WS_THICKFRAME | WS_MINIMIZEBOX | WS_MAXIMIZEBOX
	CW_USEDEFAULT       = 0x80000000
	WM_DESTROY          = 0x0002
	WM_SIZE             = 0x0005
	WM_PAINT            = 0x000F
	WM_CLOSE            = 0x0010
	WM_MOUSEMOVE        = 0x0200
	WM_LBUTTONDOWN      = 0x0201
	WM_LBUTTONUP        = 0x0202
	WM_MOUSEWHEEL       = 0x020A
	WM_KEYDOWN          = 0x0100
	WM_KEYUP            = 0x0101
	WM_CHAR             = 0x0102
	WM_TIMER            = 0x0113
	SW_SHOW             = 5
	WM_USER             = 0x0400
)

// WindowConfig defines the configuration for creating a window
type WindowConfig struct {
	Title     string
	Width     int32
	Height    int32
	X, Y      int32
	Resizable bool
}

// Window represents a GUI window
type Window struct {
	hwnd      windows.Handle
	config    WindowConfig
	EventBus  event.EventBus
	Renderer  *render.Renderer
	Root      *component.Panel
	FocusComp component.Component
}

func init() {
	// Lock OS thread for GUI operations
	runtime.LockOSThread()
}

// NewWindow creates a new window with the given configuration
func NewWindow(config WindowConfig) (*Window, error) {
	className, _ := syscall.UTF16PtrFromString("GouiWindowClass")
	title, _ := syscall.UTF16PtrFromString(config.Title)

	hInst, _, _ := procGetModuleHandleW.Call(0)

	cursor, _, _ := procLoadCursorW.Call(0, uintptr(IDC_ARROW))

	// Register Window Class
	wc := wndClassEx{
		CbSize:        uint32(unsafe.Sizeof(wndClassEx{})),
		Style:         CS_HREDRAW | CS_VREDRAW,
		LpfnWndProc:   syscall.NewCallback(wndProc),
		HInstance:     windows.Handle(hInst),
		LpszClassName: className,
		HCursor:       windows.Handle(cursor),
		HbrBackground: windows.Handle(COLOR_WINDOW + 1),
	}

	// We ignore the error if class is already registered
	procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))

	// Create Window
	style := WS_OVERLAPPEDWINDOW
	if !config.Resizable {
		style &^= WS_THICKFRAME | WS_MAXIMIZEBOX
	}

	// Default position
	x := uintptr(CW_USEDEFAULT)
	y := uintptr(CW_USEDEFAULT)
	if config.X != 0 || config.Y != 0 {
		x = uintptr(config.X)
		y = uintptr(config.Y)
	}

	hwnd, _, err := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(title)),
		uintptr(style),
		x,
		y,
		uintptr(config.Width),
		uintptr(config.Height),
		0,
		0,
		hInst,
		0,
	)

	if hwnd == 0 {
		return nil, err
	}

	// Create a timer for animations (like cursor blink)
	// ID=1, Interval=16ms (approx 60fps)
	procSetTimer.Call(uintptr(hwnd), 1, 16, 0)

	w := &Window{
		hwnd:     windows.Handle(hwnd),
		config:   config,
		EventBus: event.NewBus(),
		Renderer: render.NewRenderer(windows.Handle(hwnd), config.Width, config.Height),
		Root:     component.NewPanel(0, 0, config.Width, config.Height),
	}

	mapMu.Lock()
	windowsMap[w.hwnd] = w
	mapMu.Unlock()

	return w, nil
}

// Add adds a component to the window's root panel
func (w *Window) Add(c component.Component) {
	w.Root.Add(c)
}

// SetFocus sets the focus to a component
func (w *Window) SetFocus(c component.Component) {
	if w.FocusComp == c {
		return
	}
	if w.FocusComp != nil {
		w.FocusComp.OnBlur()
	}
	w.FocusComp = c
	if w.FocusComp != nil {
		w.FocusComp.OnFocus()
	}
}

// Show makes the window visible
func (w *Window) Show() {
	procShowWindow.Call(uintptr(w.hwnd), SW_SHOW)
	procUpdateWindow.Call(uintptr(w.hwnd))
}

// Run starts the message loop
func (w *Window) Run() {
	var msg msg
	for {
		ret, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if ret == 0 { // WM_QUIT
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}
}

// Close destroys the window
func (w *Window) Close() {
	mapMu.Lock()
	delete(windowsMap, w.hwnd)
	mapMu.Unlock()

	w.EventBus.Close()
	procDestroyWindow.Call(uintptr(w.hwnd))
}

// RequestRepaint triggers a repaint of the window.
// It is safe to call from any goroutine.
func (w *Window) RequestRepaint() {
	procPostMessageW.Call(uintptr(w.hwnd), WM_USER, 0, 0)
}

func (w *Window) Render() {
	// Optimization: check if we actually need to render?
	// For now, always render when called.

	canvas := w.Renderer.BeginFrame()
	// No need to clear if Root Panel fills everything, but safe to clear
	canvas.Clear(0xFFFFFFFF)

	w.Root.Render(canvas)

	w.Renderer.EndFrame()
	w.Renderer.Present()
}

func wndProc(hwnd windows.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	mapMu.RLock()
	w, ok := windowsMap[hwnd]
	mapMu.RUnlock()

	if ok {
		// Handle Timer for cursor blinking
		if msg == WM_TIMER {
			if w.FocusComp != nil {
				w.Render()
			}
			return 0
		}

		// Handle Repaint Request
		if msg == WM_USER {
			w.Render()
			return 0
		}

		if evt, ok := convertEvent(msg, wParam, lParam); ok {
			// Special handling for focus
			if msg == WM_LBUTTONDOWN {
				if mouseEvt, ok := evt.Data.(event.MouseEvent); ok {
					target := w.Root.FindComponentAt(mouseEvt.X, mouseEvt.Y)
					w.SetFocus(target)
				}
			}

			// Dispatch to focused component for keyboard/wheel events
			if evt.Type == event.EventKeyPress || evt.Type == event.EventKeyRelease || evt.Type == event.EventChar || evt.Type == event.EventMouseWheel {
				if w.FocusComp != nil {
					w.FocusComp.OnEvent(evt)
					w.Render() // Repaint after key event
				}
			} else {
				// Dispatch to UI components (Root) for mouse/other events
				w.Root.OnEvent(evt)
			}

			// Publish to EventBus (async)
			w.EventBus.Publish(evt)

			// Trigger repaint
			// Check if any component requested repaint (for now global repaint)
			// In a real system, we would check Root.RepaintRequested or traverse
			// For this Phase 4 optimization step, let's just say if we handled an event, we repaint.
			// But ideally we should only repaint if something changed.
			// Let's rely on the fact that components call RequestRepaint() if they change state.
			// However, we didn't implement bubble-up of repaint requests yet.
			// So we'll stick to repainting on every event for now to be safe,
			// OR we can add a flag to Window that components can set.
			w.Render()
		}

		switch msg {
		case WM_SIZE:
			width := int32(lParam & 0xFFFF)
			height := int32((lParam >> 16) & 0xFFFF)
			w.Renderer.Resize(width, height)
			w.Root.SetBounds(0, 0, width, height)
			w.Render()
		}
	}

	switch msg {
	case WM_DESTROY:
		procPostQuitMessage.Call(0)
		return 0
	}

	ret, _, _ := procDefWindowProcW.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
	return ret
}

func convertEvent(msg uint32, wParam, lParam uintptr) (event.Event, bool) {
	switch msg {
	case WM_CLOSE:
		return event.Event{Type: event.EventClose}, true
	case WM_SIZE:
		return event.Event{Type: event.EventResize}, true
	case WM_MOUSEMOVE:
		x := int32(lParam & 0xFFFF)
		y := int32((lParam >> 16) & 0xFFFF)
		return event.Event{
			Type: event.EventMouseMove,
			Data: event.MouseEvent{X: x, Y: y},
		}, true
	case WM_LBUTTONDOWN:
		x := int32(lParam & 0xFFFF)
		y := int32((lParam >> 16) & 0xFFFF)
		return event.Event{
			Type: event.EventMouseClick,
			Data: event.MouseEvent{X: x, Y: y, Button: 1},
		}, true
	case WM_LBUTTONUP:
		x := int32(lParam & 0xFFFF)
		y := int32((lParam >> 16) & 0xFFFF)
		return event.Event{
			Type: event.EventMouseRelease,
			Data: event.MouseEvent{X: x, Y: y, Button: 1},
		}, true
	case WM_MOUSEWHEEL:
		// High word of wParam is delta
		delta := int16((wParam >> 16) & 0xFFFF)
		return event.Event{
			Type: event.EventMouseWheel,
			Data: event.MouseEvent{Delta: int(delta)},
		}, true
	case WM_KEYDOWN:
		return event.Event{
			Type: event.EventKeyPress,
			Data: event.KeyEvent{
				VirtualKeyCode: uint32(wParam),
				Modifiers:      getModifiers(),
			},
		}, true
	case WM_KEYUP:
		return event.Event{
			Type: event.EventKeyRelease,
			Data: event.KeyEvent{
				VirtualKeyCode: uint32(wParam),
				Modifiers:      getModifiers(),
			},
		}, true
	case WM_CHAR:
		return event.Event{
			Type: event.EventChar,
			Data: event.KeyEvent{Rune: rune(wParam)},
		}, true
	}
	return event.Event{}, false
}

// Internal structures

type wndClassEx struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     windows.Handle
	HIcon         windows.Handle
	HCursor       windows.Handle
	HbrBackground windows.Handle
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       windows.Handle
}

type msg struct {
	Hwnd    windows.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      point
}

type point struct {
	X, Y int32
}
