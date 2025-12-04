package keys

import (
	"fmt"
	"image"
	dati "p33/DATI"
	"unsafe"

	// dati "p3/DATI"

	"github.com/fogleman/gg"
)

func Messaggio() string {
	s := "ciao"
	fmt.Println(s)
	return s
}

func ShowPos(mx, my float64) {
	// fmt.Printf("Mouse X,Y at: (%.1f, %.1f, mx:%4d, my:%4d) ---------\n", mx, my, dati.MouseX, dati.MouseY)
	dati.MouseX = int(mx)
	dati.MouseY = int(my)
}

func RGBAtoUint32Fast(img *image.RGBA) []uint32 {
	if img.Stride != img.Bounds().Dx()*4 {
		// fallback lento
		return RGBAtoUint32(img)
	}
	pix := img.Pix
	hdr := *(*[]uint32)(unsafe.Pointer(&pix))
	return hdr[:len(pix)/4]
}

func RGBAtoUint32(img *image.RGBA) []uint32 {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	out := make([]uint32, w*h)

	i := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			off := y*img.Stride + x*4
			r := uint32(img.Pix[off+0])
			g := uint32(img.Pix[off+1])
			b := uint32(img.Pix[off+2])
			a := uint32(img.Pix[off+3])

			// out[i] = (r << 24) | (g << 16) | (b << 8) | a
			out[i] = (a << 24) | (b << 16) | (g << 8) | r
			i++
		}
	}
	return out
}

// Crea overlay con una X e un cerchio
func CreateOverlay(width, height int) image.Image {
	dc := gg.NewContext(width, height)

	// sfondo trasparente
	dc.SetRGBA(0, 0, 0, 0)
	dc.Clear()

	skip := true

	if !skip {
		// X rossa
		dc.SetRGBA(1, 0, 0, 1)
		dc.SetLineWidth(10)
		dc.DrawLine(2, 2, float64(width-2), float64(height-2))
		dc.DrawLine(2, float64(height-2), float64(width-2), 2)
		dc.Stroke()

		// cerchio verde
		dc.SetRGBA(0, 1, 0, 1)
		dc.SetLineWidth(14)
		// min := math.Min(height, width)
		// min := height
		// if min > width {
		// 	min = width
		// }
		// max := height
		// if max > width {
		// 	max = width
		// }
		dc.DrawCircle(float64(width/2), float64(height/2), float64(100))
		ratioWH := float64(width) / float64(height)
		_ = ratioWH
		erx := 100.0 / (2800.0 / 1000.0) / (2400.0 / 2800.0) //ratioWH //* 0.75
		ery := 100.0                                         //* (9.0 / 16.0) /// ratioWH / ratioWH
		dc.DrawEllipse(float64(width/2), float64(height/2), erx, ery)
		dc.Stroke()

		dc.SetRGBA(0, 1, 1, .1)
		dc.DrawRectangle(10, 10, float64(width-20), float64(height-20))
		dc.Fill()

		dc.SetLineWidth(4)
		dc.SetRGBA(1, 1, 1, 1)
		dc.DrawRectangle(11, 11, float64(width-22), float64(height-22))
		dc.Stroke()

		dc.SetRGBA(1, 1, 0, 1)
		dc.DrawRectangle(0, 0, float64(width/2), float64(height/2))
		dc.Stroke()

		if err := dc.LoadFontFace("/System/Library/Fonts/SFNSRounded.ttf", 148); err != nil {
			// if err := dc.LoadFontFace("/System/Library/Fonts/Supplemental/Arial.ttf", 24); err != nil {
			panic(err)
		}
		dc.Push()
		dc.Scale(1, -1)
		dc.Translate(0, -float64(height))

		dc.SetRGBA(1, 1, 1, .99)
		rx := float64(dati.MouseX) / 1 //float64(width)
		ry := float64(dati.MouseY) / 1 // float64(height)
		// s := fmt.Sprintf("%5d %5d", dati.MouseX, dati.MouseY)
		s := fmt.Sprintf("%7.5f %7.5f  w:%5d  h:%5d", rx, ry, width, height)
		dc.DrawString(s, 100, float64(height-0))
		// dc.DrawString("CIAO", 100, 100)
		dc.Pop()

		// dc.Push()
		// dc.Scale(1, -1)
		// dc.Translate(0, -float64(height))
		// if err := dc.LoadFontFace("/System/Library/Fonts/Supplemental/Arial.ttf", 48); err != nil {
		// 	panic(err)
		// }
		// dc.DrawStringAnchored("ABCD2", float64(width)/2, float64(height)/2, 1.0, 0.5)
		// // dc.DrawString("ABCD2", float64(w)/4, float64(h)/2)
		// s := fmt.Sprintf("%5d %5d", dati.MouseX, dati.MouseY)
		// dc.DrawString(s, float64(width)/8, float64(height)/8)

		// dc.Pop()

		// // salva su file
		// dc.SavePNG("linea1.jpg")
	}
	return dc.Image()
}
