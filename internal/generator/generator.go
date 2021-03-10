package generator

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/gif"
	"io"
	"os"

	"github.com/akiomik/shiomi/internal/audio"
	simage "github.com/akiomik/shiomi/internal/image"
)

type Config struct {
	Frequency       uint
	WindowSize      uint
	SubsamplingRate uint
	LineWidth       float64
	Width           uint
	Height          uint
	BgColor         string
	FgColor         string
	Delay           uint
}

func GenerateWaveformGIF(inputFile io.Reader, outputFile io.Writer, config *Config) {
	a, err := audio.NewAudio(inputFile, config.Frequency, config.WindowSize, config.SubsamplingRate)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	colorPalette := palette.Plan9
	wimg, err := simage.NewWaveformImage(config.Width, config.Height, config.BgColor, config.FgColor, config.LineWidth)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	imgGen := make(chan *image.Paletted)
	go func() {
		defer close(imgGen)

		for audioData := range a.ReadCycles() {
			samples := audioData.Samples
			if audioData.Error != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			wimg.DrawWaveform(samples)
			pimg := wimg.ConvertToPaletted(colorPalette)
			wimg.Clear()

			imgGen <- pimg
		}
	}()

	animeGIF, err := simage.GenerateAnimationGIF(imgGen, config.Width, config.Height, colorPalette, config.Delay)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	gif.EncodeAll(outputFile, animeGIF)
}
