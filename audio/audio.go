package audio

import (
	"errors"
	"io"
	"math"

	"github.com/mjibson/go-dsp/wav"
)

type Audio struct {
	w               *wav.Wav
	freq            uint
	windowSize      uint
	subsamplingRate uint
}

type AudioData struct {
	Samples []float32
	Error   error
}

func NewAudio(f io.Reader, freq uint, windowSize uint, subsamplingRate uint) (*Audio, error) {
	if freq <= 0 {
		return nil, errors.New("freq must be > 0")
	}

	if windowSize <= 0 {
		return nil, errors.New("windowSize must be > 0")
	}

	if subsamplingRate <= 0 {
		return nil, errors.New("subsamplingRate must be > 0")
	}

	w, err := wav.New(f)
	if err != nil {
		return nil, err
	}

	audio := &Audio{w: w, freq: freq, windowSize: windowSize, subsamplingRate: subsamplingRate}
	return audio, nil
}

func (audio *Audio) NumSamplesPerCycle() float64 {
	return float64(audio.w.Header.SampleRate) / float64(audio.freq)
}

func (audio *Audio) NumSamplesPerWindow() float64 {
	return audio.NumSamplesPerCycle() * float64(audio.windowSize)
}

func (audio *Audio) RoundedNumSamplesPerWindow() uint {
	return uint(math.Round(audio.NumSamplesPerWindow()))
}

func (audio *Audio) NumWindow() uint {
	return uint(audio.w.Samples) / audio.RoundedNumSamplesPerWindow()
}

func (audio *Audio) ReadCycles() <-chan AudioData {
	nSamplesPerWindow := int(audio.RoundedNumSamplesPerWindow())
	nWindow := audio.NumWindow()

	ch := make(chan AudioData)
	go func() {
		defer close(ch)

		for i := uint(0); i < nWindow; i++ {
			samples, err := audio.w.ReadFloats(nSamplesPerWindow)
			if i%audio.subsamplingRate != 0 {
				continue
			}

			ch <- AudioData{Samples: samples, Error: err}
		}
	}()

	return ch
}
