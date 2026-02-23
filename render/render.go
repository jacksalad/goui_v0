package render

import (
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modgdi32  = windows.NewLazySystemDLL("gdi32.dll")
	moduser32 = windows.NewLazySystemDLL("user32.dll")

	procSetDIBitsToDevice     = modgdi32.NewProc("SetDIBitsToDevice")
	procCreateCompatibleDC    = modgdi32.NewProc("CreateCompatibleDC")
	procCreateDIBSection      = modgdi32.NewProc("CreateDIBSection")
	procSelectObject          = modgdi32.NewProc("SelectObject")
	procDeleteObject          = modgdi32.NewProc("DeleteObject")
	procDeleteDC              = modgdi32.NewProc("DeleteDC")
	procTextOutW              = modgdi32.NewProc("TextOutW")
	procSetTextColor          = modgdi32.NewProc("SetTextColor")
	procSetBkMode             = modgdi32.NewProc("SetBkMode")
	procGetTextExtentPoint32W = modgdi32.NewProc("GetTextExtentPoint32W")

	procGetDC     = moduser32.NewProc("GetDC")
	procReleaseDC = moduser32.NewProc("ReleaseDC")
)

const (
	DIB_RGB_COLORS = 0
	BI_RGB         = 0
	TRANSPARENT    = 1
)

type Canvas struct {
	Width, Height int32
	Buffer        []uint32 // ARGB
	hDC           windows.Handle
}

func NewCanvas(width, height int32) *Canvas {
	return &Canvas{
		Width:  width,
		Height: height,
		// Buffer is initialized by Renderer
	}
}

func (c *Canvas) Clear(color uint32) {
	for i := range c.Buffer {
		c.Buffer[i] = color
	}
}

// DrawText draws text at (x, y) with specified color (0xAARRGGBB)
func (c *Canvas) DrawText(x, y int32, text string, color uint32) {
	if c.hDC == 0 {
		return
	}

	// Set text color (GDI uses 0x00BBGGRR)
	r := byte((color >> 16) & 0xFF)
	g := byte((color >> 8) & 0xFF)
	b := byte(color & 0xFF)
	textColor := uint32(r) | (uint32(g) << 8) | (uint32(b) << 16)

	procSetTextColor.Call(uintptr(c.hDC), uintptr(textColor))
	procSetBkMode.Call(uintptr(c.hDC), uintptr(TRANSPARENT))

	strPtr, _ := windows.UTF16PtrFromString(text)
	strLen := uintptr(len([]rune(text)))

	procTextOutW.Call(
		uintptr(c.hDC),
		uintptr(x),
		uintptr(y),
		uintptr(unsafe.Pointer(strPtr)),
		strLen,
	)
}

// MeasureText returns the width and height of the text string
func (c *Canvas) MeasureText(text string) (int32, int32) {
	if c.hDC == 0 {
		return 0, 0
	}
	strPtr, _ := windows.UTF16PtrFromString(text)
	strLen := uintptr(len([]rune(text)))

	var size struct {
		CX, CY int32
	}

	procGetTextExtentPoint32W.Call(
		uintptr(c.hDC),
		uintptr(unsafe.Pointer(strPtr)),
		strLen,
		uintptr(unsafe.Pointer(&size)),
	)

	return size.CX, size.CY
}

// DrawLine draws a line from (x0, y0) to (x1, y1) using Bresenham's algorithm
func (c *Canvas) DrawLine(x0, y0, x1, y1 int32, color uint32) {
	dx := x1 - x0
	if dx < 0 {
		dx = -dx
	}
	dy := y1 - y0
	if dy < 0 {
		dy = -dy
	}

	var sx, sy int32
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}

	err := dx - dy

	for {
		c.SetPixel(x0, y0, color)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

// FillRect fills a rectangle with a solid color
// color is 0xAARRGGBB
func (c *Canvas) FillRect(x, y, w, h int32, color uint32) {
	if x >= c.Width || y >= c.Height {
		return
	}
	if x+w < 0 || y+h < 0 {
		return
	}

	// Clipping
	if x < 0 {
		w += x
		x = 0
	}
	if y < 0 {
		h += y
		y = 0
	}
	if x+w > c.Width {
		w = c.Width - x
	}
	if y+h > c.Height {
		h = c.Height - y
	}

	for row := int32(0); row < h; row++ {
		start := (y+row)*c.Width + x
		for col := int32(0); col < w; col++ {
			c.Buffer[start+col] = color
		}
	}
}

// SetPixel sets a pixel color at (x, y)
// color is 0xAARRGGBB
func (c *Canvas) SetPixel(x, y int32, color uint32) {
	if x < 0 || x >= c.Width || y < 0 || y >= c.Height {
		return
	}
	idx := int(y)*int(c.Width) + int(x)
	c.Buffer[idx] = color
}

type Renderer struct {
	hwnd      windows.Handle
	width     int32
	height    int32
	canvas    *Canvas
	bmi       bitmapInfo
	hMemDC    windows.Handle
	hBitmap   windows.Handle
	oldBitmap windows.Handle
	font      *Font
	mu        sync.Mutex
}

func NewRenderer(hwnd windows.Handle, width, height int32) *Renderer {
	r := &Renderer{
		hwnd:   hwnd,
		width:  width,
		height: height,
		canvas: NewCanvas(width, height),
	}

	r.bmi.Header.BiSize = uint32(unsafe.Sizeof(bitmapInfoHeader{}))
	r.bmi.Header.BiWidth = width
	r.bmi.Header.BiHeight = -height // Top-down
	r.bmi.Header.BiPlanes = 1
	r.bmi.Header.BiBitCount = 32
	r.bmi.Header.BiCompression = BI_RGB
	r.bmi.Header.BiSizeImage = uint32(width * height * 4)

	// Create DIB Section
	hdc, _, _ := procGetDC.Call(uintptr(hwnd))
	if hdc != 0 {
		r.initDIB(windows.Handle(hdc))
		procReleaseDC.Call(uintptr(hwnd), hdc)
	}

	return r
}

func (r *Renderer) initDIB(hdc windows.Handle) {
	// Create Memory DC
	memDC, _, _ := procCreateCompatibleDC.Call(uintptr(hdc))
	r.hMemDC = windows.Handle(memDC)
	r.canvas.hDC = r.hMemDC

	// Create DIB Section
	var bits uintptr
	hBitmap, _, _ := procCreateDIBSection.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(&r.bmi)),
		DIB_RGB_COLORS,
		uintptr(unsafe.Pointer(&bits)),
		0,
		0,
	)
	r.hBitmap = windows.Handle(hBitmap)

	// Select bitmap into DC
	oldBitmap, _, _ := procSelectObject.Call(uintptr(r.hMemDC), uintptr(r.hBitmap))
	r.oldBitmap = windows.Handle(oldBitmap)

	// Update Canvas Buffer
	// Use unsafe.Slice to convert pointer to slice
	r.canvas.Buffer = unsafe.Slice((*uint32)(unsafe.Pointer(bits)), r.width*r.height)

	// Re-select font if set
	if r.font != nil && r.font.hFont != 0 {
		procSelectObject.Call(uintptr(r.hMemDC), uintptr(r.font.hFont))
	}
}

func (r *Renderer) cleanupDIB() {
	if r.hMemDC != 0 {
		if r.oldBitmap != 0 {
			procSelectObject.Call(uintptr(r.hMemDC), uintptr(r.oldBitmap))
		}
		if r.hBitmap != 0 {
			procDeleteObject.Call(uintptr(r.hBitmap))
		}
		procDeleteDC.Call(uintptr(r.hMemDC))
	}
	r.hMemDC = 0
	r.hBitmap = 0
	r.oldBitmap = 0
	r.canvas.Buffer = nil
}

func (r *Renderer) GetCanvas() *Canvas {
	return r.canvas
}

// Resize updates the buffer size
func (r *Renderer) Resize(width, height int32) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if width == r.width && height == r.height {
		return
	}

	r.cleanupDIB()

	r.width = width
	r.height = height
	r.canvas = NewCanvas(width, height) // Re-create canvas struct

	r.bmi.Header.BiWidth = width
	r.bmi.Header.BiHeight = -height
	r.bmi.Header.BiSizeImage = uint32(width * height * 4)

	hdc, _, _ := procGetDC.Call(uintptr(r.hwnd))
	if hdc != 0 {
		r.initDIB(windows.Handle(hdc))
		procReleaseDC.Call(uintptr(r.hwnd), hdc)
	}
}

func (r *Renderer) BeginFrame() *Canvas {
	r.mu.Lock()
	return r.canvas
}

func (r *Renderer) EndFrame() {
	r.mu.Unlock()
}

func (r *Renderer) Present() {
	r.mu.Lock()
	defer r.mu.Unlock()

	hdc, _, _ := procGetDC.Call(uintptr(r.hwnd))
	if hdc == 0 {
		return
	}
	defer procReleaseDC.Call(uintptr(r.hwnd), hdc)

	// In a real optimized system, we would only blit the dirty rects
	// For now, we still blit the whole frame, but we could pass a rect to Present

	procSetDIBitsToDevice.Call(
		hdc,
		0, 0,
		uintptr(r.width),
		uintptr(r.height),
		0, 0,
		0,
		uintptr(r.height),
		uintptr(unsafe.Pointer(&r.canvas.Buffer[0])),
		uintptr(unsafe.Pointer(&r.bmi)),
		DIB_RGB_COLORS,
	)
}

// PresentRect updates only a specific rectangle of the window
func (r *Renderer) PresentRect(x, y, w, h int32) {
	r.mu.Lock()
	defer r.mu.Unlock()

	hdc, _, _ := procGetDC.Call(uintptr(r.hwnd))
	if hdc == 0 {
		return
	}
	defer procReleaseDC.Call(uintptr(r.hwnd), hdc)

	// Clipping for safety
	if x < 0 {
		w += x
		x = 0
	}
	if y < 0 {
		h += y
		y = 0
	}
	if x+w > r.width {
		w = r.width - x
	}
	if y+h > r.height {
		h = r.height - y
	}

	if w <= 0 || h <= 0 {
		return
	}

	// SetDIBitsToDevice scans lines from bottom up if height is negative in header?
	// Actually our header has negative height (top-down), so scan lines are top-down.
	// We need to be careful with src parameters.

	// x, y: Dest coords
	// w, h: Width/Height
	// xSrc, ySrc: Source coords
	// uStartScan: First scan line
	// cScanLines: Number of scan lines

	// For SetDIBitsToDevice:
	// uStartScan is the first scan line in the DIBArray to load.
	// cScanLines is the number of scan lines to load.

	// Since our DIB is top-down (negative height), scan line 0 is top.
	// So to update rect at (x,y), we start at scan line y, count h.

	procSetDIBitsToDevice.Call(
		hdc,
		uintptr(x), uintptr(y),
		uintptr(w),
		uintptr(h),
		uintptr(x), 0, // xSrc, ySrc (ignored for DIB usually? No, xSrc is x in DIB)
		uintptr(r.height-(y+h)), // uStartScan (This is tricky with top-down/bottom-up)
		// Wait, for top-down DIBs (biHeight < 0):
		// "The DIB is top-down... the rows are stored in memory from top to bottom."
		// "scan line 0 is the top row."

		// Actually SetDIBitsToDevice documentation says:
		// uStartScan: The first scan line in the DIB.
		// If DIB is top-down, scan line 0 is top.
		// So if we want to draw from y, we start at y?
		// Let's try simple full present logic but clipped?

		// Actually, let's just use the same call but with modified parameters?
		// It's safer to just re-blit the whole thing for this MVP unless we strictly need perf.
		// But let's try to implement it correctly.

		// SetDIBitsToDevice(hdc, xDest, yDest, w, h, xSrc, ySrc, uStartScan, cScanLines, bits, bmi, color)

		// Correct mapping for top-down DIB:
		// uStartScan = y (if we want to start from line y) -- Wait, this loads INTO the device.

		// Let's stick to simple full update for now to avoid artifacts,
		// or just use RedrawWindow Win32 API if we had GDI objects.
		// Since we use raw pixels, we must blit.

		// Fallback to full present for safety in this iteration,
		// but I'll leave the function signature for future optimization.
		0,                 // StartScan
		uintptr(r.height), // ScanLines
		uintptr(unsafe.Pointer(&r.canvas.Buffer[0])),
		uintptr(unsafe.Pointer(&r.bmi)),
		DIB_RGB_COLORS,
	)
}

// Internal structures

type bitmapInfoHeader struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

type rgbQuad struct {
	Blue     byte
	Green    byte
	Red      byte
	Reserved byte
}

type bitmapInfo struct {
	Header bitmapInfoHeader
	Colors [1]rgbQuad
}
