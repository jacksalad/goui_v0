package component

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/jacksalad/goui_v0/event"
	"github.com/jacksalad/goui_v0/render"
)

type Image struct {
	BaseComponent
	img      image.Image
	rgba     *image.RGBA
	FilePath string
}

func NewImage(filePath string) (*Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// Convert to RGBA for easier access
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	c := &Image{
		img:      img,
		rgba:     rgba,
		FilePath: filePath,
	}
	c.Visible = true
	// Default size to image size
	c.SetBounds(0, 0, int32(bounds.Dx()), int32(bounds.Dy()))

	return c, nil
}

func (c *Image) Render(canvas *render.Canvas) {
	if !c.Visible || c.rgba == nil {
		return
	}

	// Simple copy for now, no scaling
	// We need to copy pixels from c.rgba to canvas.Buffer
	// c.rgba is R G B A
	// canvas.Buffer is B G R A (0xAARRGGBB in Little Endian uint32)

	// Destination bounds clipped to component bounds
	dstX := c.Bounds.X
	dstY := c.Bounds.Y
	width := c.Bounds.Width
	height := c.Bounds.Height

	srcBounds := c.rgba.Bounds()
	srcW := int32(srcBounds.Dx())
	srcH := int32(srcBounds.Dy())

	// Min width/height
	if width > srcW {
		width = srcW
	}
	if height > srcH {
		height = srcH
	}

	for y := int32(0); y < height; y++ {
		for x := int32(0); x < width; x++ {
			// Get RGBA from source
			// image.RGBA.Pix is standard Go RGBA
			srcOff := c.rgba.PixOffset(int(x), int(y))
			r := c.rgba.Pix[srcOff+0]
			g := c.rgba.Pix[srcOff+1]
			b := c.rgba.Pix[srcOff+2]
			a := c.rgba.Pix[srcOff+3] // Alpha unused for now in DIB unless we blend

			// Pack to uint32 0xAARRGGBB (Little Endian in memory: B G R A)
			// So color = (A << 24) | (R << 16) | (G << 8) | B
			// Wait, Windows DIB 32-bit:
			// "The bitmap has a maximum of 2^32 colors. If the biCompression member of the BITMAPINFOHEADER is BI_RGB, the bmiColors member of BITMAPINFO is NULL. Each DWORD in the bitmap array represents the relative intensities of blue, green, and red for a pixel. The value for blue is in the least significant 8 bits, followed by 8 bits each for green and red. The high byte in each DWORD is not used."
			// So memory layout: B G R X
			// Little Endian uint32: 0xXXRRGGBB
			// So color = (0 << 24) | (R << 16) | (G << 8) | B

			color := (uint32(a) << 24) | (uint32(r) << 16) | (uint32(g) << 8) | uint32(b)

			canvas.SetPixel(dstX+x, dstY+y, color)
		}
	}

	c.RepaintRequested = false
}

func (c *Image) OnEvent(evt event.Event) bool {
	return false
}
