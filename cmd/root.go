package cmd

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/gif"
	"io"
	"os"

	"github.com/akiomik/shiomi/config"
	"github.com/akiomik/shiomi/internal/audio"
	simage "github.com/akiomik/shiomi/internal/image"
	"github.com/spf13/cobra"
)

var (
	input      string
	output     string
	freq       uint
	windowSize uint
	rate       uint
	width      uint
	height     uint
	bgColor    string
	fgColor    string
)

var rootCmd = &cobra.Command{
	Use:     "shiomi",
	Short:   "shiomi " + config.Version,
	Version: config.Version,
	Run: func(cmd *cobra.Command, args []string) {
		if freq <= 0 {
			fmt.Println("freq must be > 0")
			os.Exit(1)
		}

		if windowSize <= 0 {
			fmt.Println("num must be > 0")
			os.Exit(1)
		}

		if rate <= 0 {
			fmt.Println("rate must be > 0")
			os.Exit(1)
		}

		if width <= 0 || height <= 0 {
			fmt.Println("width and height must be > 0")
			os.Exit(1)
		}

		inputFile, err := os.Open(input)
		defer inputFile.Close()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		outputFile, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		defer outputFile.Close()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		generateWaveformGif(inputFile, outputFile)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&input, "input", "i", "", "A *.wav file (required)")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "An output gif file (required)")
	rootCmd.Flags().UintVarP(&freq, "freq", "f", 1000, "The frequency of an input file [Hz]")
	rootCmd.Flags().UintVarP(&windowSize, "num", "n", 3, "An output number of cycles")
	rootCmd.Flags().UintVarP(&rate, "rate", "r", 10, "A subsampling rate for output")
	rootCmd.Flags().UintVarP(&width, "width", "", 320, "A width of an output file [px]")
	rootCmd.Flags().UintVarP(&height, "height", "", 240, "A height of an output file [px]")
	rootCmd.Flags().StringVarP(&bgColor, "bg", "", "#000000", "A background color for output")
	rootCmd.Flags().StringVarP(&fgColor, "fg", "", "#ffffff", "A foreground color for output")
	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func generateWaveformGif(inputFile io.Reader, outputFile io.Writer) {
	a, err := audio.NewAudio(inputFile, freq, windowSize, rate)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	colorPalette := palette.Plan9
	lineWidth := 2.0
	wimg, err := simage.NewWaveformImage(width, height, bgColor, fgColor, lineWidth)
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

	animeGif, err := simage.GenerateAnimationGif(imgGen, width, height, colorPalette, 0)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	gif.EncodeAll(outputFile, animeGif)
}
