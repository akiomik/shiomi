package image

import (
	"errors"
	"image"
	"image/color"
	"image/gif"
)

func GenerateAnimationGIF(ch <-chan *image.Paletted, width uint, height uint, palette color.Palette, delay uint) (*gif.GIF, error) {
	if width <= 0 || height <= 0 {
		return nil, errors.New("width and height must be > 0")
	}

	var images []*image.Paletted
	var delays []int

	for img := range ch {
		images = append(images, img)
		delays = append(delays, int(delay))
	}

	colorModel := image.NewPaletted(image.Rect(0, 0, 1, 1), palette).ColorModel()
	gifConfig := image.Config{ColorModel: colorModel, Width: int(width), Height: int(height)}
	return &gif.GIF{
		Image:  images,
		Delay:  delays,
		Config: gifConfig,
	}, nil
}
