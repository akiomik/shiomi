package cmd

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/gif"
	"math"
	"os"

	"github.com/akiomik/shiomi/audio"
	"github.com/akiomik/shiomi/config"
	"github.com/fogleman/gg"
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

		a, err := audio.NewAudio(inputFile, freq, windowSize, rate)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		var images []*image.Paletted
		var delays []int

		bgColor := color.RGBA{0x00, 0x00, 0x00, 0xff}
		fgColor := color.RGBA{0x00, 0xff, 0x00, 0xff}
		colorPalette := palette.Plan9

		heightf := float64(height)
		widthf := float64(width)
		xScale := widthf / a.NumSamplesPerWindow()

		readCycles := a.ReadCycles()
		for audioData := range readCycles {
			samples := audioData.Samples
			if audioData.Error != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			dc := gg.NewContext(width, height)
			dc.SetColor(bgColor)
			dc.DrawRectangle(0, 0, widthf, heightf)
			dc.Fill()
			dc.SetColor(fgColor)
			dc.SetLineWidth(0.8)

			for j, sample := range samples[:len(samples)-1] {
				x1 := math.Round(float64(j) * xScale)
				y1 := heightf - math.Round(float64(sample)*heightf)
				x2 := math.Round(float64(j+1) * xScale)
				y2 := heightf - math.Round(float64(audioData.Samples[j+1])*heightf)
				dc.DrawLine(x1, y1, x2, y2)
			}

			dc.Stroke()

			img := image.NewPaletted(dc.Image().Bounds(), colorPalette)
			bounds := img.Bounds()
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
					img.Set(x, y, dc.Image().At(x, y))
				}
			}

			images = append(images, img)
			delays = append(delays, 0)
		}

		outputFile, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0600)
		defer outputFile.Close()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		colorModel := image.NewPaletted(image.Rect(0, 0, 1, 1), colorPalette).ColorModel()
		gifConfig := image.Config{ColorModel: colorModel, Width: width, Height: height}
		gif.EncodeAll(outputFile, &gif.GIF{
			Image:  images,
			Delay:  delays,
			Config: gifConfig,
		})
	},
}

func init() {
	rootCmd.Flags().StringVarP(&input, "input", "i", "", "A *.wav file (required)")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "An output gif file (required)")
	rootCmd.Flags().UintVarP(&freq, "freq", "f", 1000, "The frequency of an input file [Hz]")
	rootCmd.Flags().UintVarP(&windowSize, "num", "n", 3, "An output number of cycles")
	rootCmd.Flags().UintVarP(&rate, "rate", "r", 10, "A subsampling rate for output")
	rootCmd.Flags().IntVarP(&width, "width", "", 640, "A width of an output file")
	rootCmd.Flags().IntVarP(&height, "height", "", 480, "A height of an output file")
	rootCmd.MarkFlagRequired("input")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
