package cmd

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/gif"
	"io"
	"os"

	"github.com/akiomik/shiomi/audio"
	"github.com/akiomik/shiomi/config"
	simage "github.com/akiomik/shiomi/image"
	"github.com/spf13/cobra"
)

var (
	input      string
	output     string
	freq       uint
	windowSize uint
	rate       uint
	width      int
	height     int
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

		outputFile, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0600)
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
	rootCmd.Flags().IntVarP(&width, "width", "", 640, "A width of an output file [px]")
	rootCmd.Flags().IntVarP(&height, "height", "", 480, "A height of an output file [px]")
	rootCmd.Flags().StringVarP(&bgColor, "bg", "", "#000000", "A background color for output")
	rootCmd.Flags().StringVarP(&fgColor, "fg", "", "#00ff00", "A foreground color for output")
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

	var images []*image.Paletted
	var delays []int
	colorPalette := palette.Plan9

	wimg, err := simage.NewWaveformImage(uint(width), uint(height), bgColor, fgColor, 0.8)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	readCycles := a.ReadCycles()
	for audioData := range readCycles {
		samples := audioData.Samples
		if audioData.Error != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		wimg.DrawWaveform(samples)
		pimg := wimg.ConvertToPaletted(colorPalette)
		wimg.Clear()

		images = append(images, pimg)
		delays = append(delays, 0)
	}

	colorModel := image.NewPaletted(image.Rect(0, 0, 1, 1), colorPalette).ColorModel()
	gifConfig := image.Config{ColorModel: colorModel, Width: width, Height: height}
	gif.EncodeAll(outputFile, &gif.GIF{
		Image:  images,
		Delay:  delays,
		Config: gifConfig,
	})
}
