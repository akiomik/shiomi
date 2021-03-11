package image

import (
	"image/color"
	"image/color/palette"
	"testing"
)

func TestConvertImageToPaletted(t *testing.T) {
	wimg, err := NewWaveformImage(128, 64, "#ff00ff", "#00ff00", 1)
	if wimg == nil {
		t.Errorf(err.Error())
	}

	wimg.DrawBackground()
	pimg := ConvertImageToPaletted(wimg.Image(), palette.Plan9)
	bgColor := color.RGBA{0xff, 0x00, 0xff, 0xff}

	for x := uint(0); x < wimg.Width; x++ {
		for y := uint(0); y < wimg.Height; y++ {
			got := pimg.At(int(x), int(y))
			if got != bgColor {
				r, g, b, a := got.RGBA()
				t.Errorf("ConvertImageToPaletted(img, palette.Plan9) draws (0x%02x, 0x%02x, 0x%02x, 0x%02x) at (%d, %d); want (0xff, 0x00, 0xff, 0xff)", r, g, b, a, x, y)
			}
		}
	}
}
