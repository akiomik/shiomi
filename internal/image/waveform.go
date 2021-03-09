package image

import (
	"errors"
	"image"
	"image/color"
	"math"

	"github.com/fogleman/gg"
)

type WaveformImage struct {
	Width     uint
	Height    uint
	BgColor   string
	FgColor   string
	LineWidth float64
	context   *gg.Context
}

type Point = gg.Point

func NewWaveformImage(width uint, height uint, bgColor string, fgColor string, lineWidth float64) (*WaveformImage, error) {
	if width <= 0 || height <= 0 {
		return nil, errors.New("width and height must be > 0")
	}

	context := gg.NewContext(int(width), int(height))
	return &WaveformImage{Width: width, Height: height, BgColor: bgColor, FgColor: fgColor, LineWidth: lineWidth, context: context}, nil
}

func (wimg *WaveformImage) DrawBackground() {
	wimg.context.SetHexColor(wimg.BgColor)
	wimg.context.DrawRectangle(0, 0, float64(wimg.Width), float64(wimg.Height))
	wimg.context.Fill()
}

func (wimg *WaveformImage) DrawPolyLine(points []Point) {
	wimg.context.SetHexColor(wimg.FgColor)
	wimg.context.SetLineWidth(wimg.LineWidth)

	for i, p1 := range points[:len(points)-1] {
		p2 := points[i+1]
		wimg.context.DrawLine(p1.X, p1.Y, p2.X, p2.Y)
	}

	wimg.context.Stroke()
}

func (wimg *WaveformImage) DrawWaveform(data []float32) {
	widthf := float64(wimg.Width)
	heightf := float64(wimg.Height)
	xScale := widthf / float64(len(data)-1)

	wimg.DrawBackground()

	var points = make([]gg.Point, len(data))
	for i, d := range data {
		x := math.Round(float64(i) * xScale)
		y := heightf - math.Round(float64(d)*heightf)
		points[i] = gg.Point{X: x, Y: y}
	}

	wimg.DrawPolyLine(points)
}

func (wimg *WaveformImage) Clear() {
	wimg.context.Clear()
}

func (wimg *WaveformImage) Image() image.Image {
	return wimg.context.Image()
}

func (wimg *WaveformImage) ConvertToPaletted(colorPalette color.Palette) *image.Paletted {
	img := wimg.Image()
	bounds := img.Bounds()
	pimg := image.NewPaletted(bounds, colorPalette)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			pimg.Set(x, y, img.At(x, y))
		}
	}

	return pimg
}
