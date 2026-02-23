package render

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	procAddFontResourceExW    = modgdi32.NewProc("AddFontResourceExW")
	procRemoveFontResourceExW = modgdi32.NewProc("RemoveFontResourceExW")
	procCreateFontW           = modgdi32.NewProc("CreateFontW")
	procGetDeviceCaps         = modgdi32.NewProc("GetDeviceCaps")
)

const (
	FR_PRIVATE          = 0x10
	FR_NOT_ENUM         = 0x20
	LOGPIXELSY          = 90
	FW_NORMAL           = 400
	FW_BOLD             = 700
	DEFAULT_CHARSET     = 1
	OUT_DEFAULT_PRECIS  = 0
	CLIP_DEFAULT_PRECIS = 0
	DEFAULT_QUALITY     = 0
	DEFAULT_PITCH       = 0
	FF_DONTCARE         = 0
)

type Font struct {
	hFont windows.Handle
	Name  string
	Size  int
}

// LoadFontFile loads a font from a file path.
// The font is private to the application.
// You still need to know the font face name to create a Font object.
func LoadFontFile(path string) error {
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	ret, _, _ := procAddFontResourceExW.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		FR_PRIVATE,
		0,
	)
	if ret == 0 {
		// Try to get error
		err := syscall.GetLastError()
		if err != nil {
			return err
		}
		// If no error code but returned 0, it means 0 fonts added.
		// This can happen if the file is invalid.
		return syscall.EINVAL
	}
	return nil
}

// NewFont creates a new font with the given family name and size (in points).
// If the font is not found, Windows will substitute it.
func NewFont(name string, size int) *Font {
	// We need a DC to calculate height, use screen DC
	hdc, _, _ := procGetDC.Call(0)
	defer procReleaseDC.Call(0, hdc)

	logPixelsY, _, _ := procGetDeviceCaps.Call(hdc, LOGPIXELSY)
	height := -int32(size * int(logPixelsY) / 72)

	namePtr, _ := syscall.UTF16PtrFromString(name)

	hFont, _, _ := procCreateFontW.Call(
		uintptr(height),
		0,                         // Width
		0,                         // Escapement
		0,                         // Orientation
		FW_NORMAL,                 // Weight
		0,                         // Italic
		0,                         // Underline
		0,                         // StrikeOut
		DEFAULT_CHARSET,           // CharSet
		OUT_DEFAULT_PRECIS,        // OutPrecision
		CLIP_DEFAULT_PRECIS,       // ClipPrecision
		DEFAULT_QUALITY,           // Quality
		DEFAULT_PITCH|FF_DONTCARE, // PitchAndFamily
		uintptr(unsafe.Pointer(namePtr)),
	)

	return &Font{
		hFont: windows.Handle(hFont),
		Name:  name,
		Size:  size,
	}
}

// Close destroys the font object.
func (f *Font) Close() {
	if f.hFont != 0 {
		procDeleteObject.Call(uintptr(f.hFont))
		f.hFont = 0
	}
}

// SetFont sets the current font for the renderer.
// This font will be used for all subsequent text drawing operations.
func (r *Renderer) SetFont(font *Font) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.font = font
	if r.hMemDC != 0 && font != nil && font.hFont != 0 {
		procSelectObject.Call(uintptr(r.hMemDC), uintptr(font.hFont))
	}
}

// SetFont sets the current font for the canvas.
// Note: This only sets it for the current frame/DC.
// Use Renderer.SetFont to persist across resize.
func (c *Canvas) SetFont(font *Font) {
	if c.hDC == 0 || font == nil || font.hFont == 0 {
		return
	}
	procSelectObject.Call(uintptr(c.hDC), uintptr(font.hFont))
}

// MeasureText calculates the width and height of the given text with the specified font.
// If font is nil, it uses the system default font.
func MeasureText(text string, font *Font) (int32, int32) {
	hdc, _, _ := procGetDC.Call(0)
	defer procReleaseDC.Call(0, hdc)

	oldFont := uintptr(0)
	if font != nil && font.hFont != 0 {
		oldFont, _, _ = procSelectObject.Call(hdc, uintptr(font.hFont))
	}

	utf16Str, _ := syscall.UTF16FromString(text)
	var size struct {
		CX, CY int32
	}
	procGetTextExtentPoint32W.Call(
		hdc,
		uintptr(unsafe.Pointer(&utf16Str[0])),
		uintptr(len(utf16Str)-1), // -1 for null terminator? No, syscall adds it but len includes it?
		// syscall.UTF16FromString returns []uint16 with null terminator.
		// GetTextExtentPoint32W expects length of string.
		// So we should use len(utf16Str) - 1.
		// Wait, let's verify. Yes.
		uintptr(unsafe.Pointer(&size)),
	)

	if oldFont != 0 {
		procSelectObject.Call(hdc, oldFont)
	}

	return size.CX, size.CY
}
