package image

import (
	"image/color"
	"testing"
)

func TestNewWaveformImage(t *testing.T) {
	_, err := NewWaveformImage(0, 64, "#000000", "#00ff00", 1)
	if err == nil {
		t.Errorf("NewWaveformImage(0, 64, \"#000000\", \"#00ff00\", 1) = _, nil; want error")
	}

	_, err = NewWaveformImage(128, 0, "#000000", "#00ff00", 1)
	if err == nil {
		t.Errorf("NewWaveformImage(128, 0, \"#000000\", \"#00ff00\", 1) = _, nil; want error")
	}

	wimg, _ := NewWaveformImage(128, 64, "#000000", "#00ff00", 1)
	if wimg == nil {
		t.Errorf("NewWaveformImage(128, 64, \"#000000\", \"#00ff00\", 1) = nil, _; want WaveformImage")
	}

	got := wimg.context.Width()
	if got != 128 {
		t.Errorf("NewWaveformImage(128, 64, \"#000000\", \"#00ff00\", 1) = WaveformImage{context: {Width: %d}}, _; want WaveformImage{context: {Width: 128}}", got)
	}

	got = wimg.context.Height()
	if got != 64 {
		t.Errorf("NewWaveformImage(128, 64, \"#000000\", \"#00ff00\", 1) = WaveformImage{context: {Height: %d}}, _; want WaveformImage{context: {Height: 64}}", got)
	}
}

func TestDrawBackground(t *testing.T) {
	wimg, err := NewWaveformImage(128, 64, "#000000", "#00ff00", 1)
	if wimg == nil {
		t.Errorf(err.Error())
	}

	wimg.DrawBackground()
	img := wimg.Image()
	bgColor := color.RGBA{0x00, 0x00, 0x00, 0xff}

	for x := uint(0); x < wimg.Width; x++ {
		for y := uint(0); y < wimg.Height; y++ {
			got := img.At(int(x), int(y))
			if got != bgColor {
				r, g, b, a := got.RGBA()
				t.Errorf("DrawBackground() draws (0x%02x, 0x%02x, 0x%02x, 0x%02x) at (%d, %d); want (0x00, 0x00, 0x00, 0xff)", r, g, b, a, x, y)
			}
		}
	}
}

func TestDrawPolyLine(t *testing.T) {
	// wimg, err := NewWaveformImage(128, 64, "#000000", "#00ff00", 1)
	// if wimg == nil {
	// 	t.Errorf(err.Error())
	// }
	//
	// left, top := 0.0, 0.0
	// right, bottom := float64(wimg.Width)-1, float64(wimg.Height)-1
	//
	// img := wimg.Image()
	// bgColor := color.RGBA{0x00, 0x00, 0x00, 0x00}
	//
	// for x := uint(0); x < wimg.Width; x++ {
	// 	for y := uint(0); y < wimg.Height; y++ {
	// 		got := img.At(int(x), int(y))
	// 		xf, yf := float64(x), float64(y)
	//
	// 		if (xf == left || xf == right || yf == top || yf == bottom) && got == bgColor {
	// 			r, g, b, a := got.RGBA()
	// 			t.Errorf("DrawPolyLine() draws (0x%02x, 0x%02x, 0x%02x, 0x%02x) at (%d, %d); want (0x00, 0xff, 0x00, 0xff)", r, g, b, a, x, y)
	// 		} else if xf != left && xf != right && yf != top && yf != bottom && got != bgColor {
	// 			r, g, b, a := got.RGBA()
	// 			t.Errorf("DrawPolyLine() draws (0x%02x, 0x%02x, 0x%02x, 0x%02x) at (%d, %d); want (0x00, 0x00, 0x00, 0x00)", r, g, b, a, x, y)
	// 		}
	// 	}
	// }
}
