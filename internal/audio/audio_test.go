package audio

import (
	"os"
	"testing"
)

var testTone = "../../test/data/sine-44.1kHz-16bit-1kHz-5s.wav"

func TestNewAudio(t *testing.T) {
	f, err := os.Open(testTone)
	defer f.Close()
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = NewAudio(f, 0, 3, 10)
	if err == nil {
		t.Errorf("NewAudio(f, 0, 3, 10) = _, nil; want error")
	}

	_, err = NewAudio(f, 1000, 0, 10)
	if err == nil {
		t.Errorf("NewAudio(f, 1000, 0, 10) = _, nil; want error")
	}

	_, err = NewAudio(f, 1000, 3, 0)
	if err == nil {
		t.Errorf("NewAudio(f, 1000, 3, 0) = _, nil; want error")
	}

	a, err := NewAudio(f, 1000, 3, 10)
	if a == nil {
		t.Errorf("NewAudio(f, 1000, 3, 10) = nil, _; want Audio")
	}
}

func TestNumSamplesPerCycle(t *testing.T) {
	f, err := os.Open(testTone)
	defer f.Close()
	if err != nil {
		t.Errorf(err.Error())
	}

	a, err := NewAudio(f, 1000, 3, 10)
	if err != nil {
		t.Errorf(err.Error())
	}

	got := a.NumSamplesPerCycle()
	if got != 44.1 {
		t.Errorf("NumSamplesPerCycle() = %f; want 44.1", got)
	}
}

func TestNumSamplesPerWindow(t *testing.T) {
	f, err := os.Open(testTone)
	defer f.Close()
	if err != nil {
		t.Errorf(err.Error())
	}

	a, err := NewAudio(f, 1000, 3, 10)
	if err != nil {
		t.Errorf(err.Error())
	}

	got := a.NumSamplesPerWindow()
	if got != 132.3 {
		t.Errorf("NumSamplesPerWindow() = %f; want 132.3", got)
	}
}

func TestRoundedNumSamplesPerWindow(t *testing.T) {
	f, err := os.Open(testTone)
	defer f.Close()
	if err != nil {
		t.Errorf(err.Error())
	}

	a, err := NewAudio(f, 1000, 3, 10)
	if err != nil {
		t.Errorf(err.Error())
	}

	got := a.RoundedNumSamplesPerWindow()
	if got != 132 {
		t.Errorf("RoundedNumSamplesPerWindow() = %d; want 132", got)
	}
}

func TestNumWindow(t *testing.T) {
	f, err := os.Open(testTone)
	defer f.Close()
	if err != nil {
		t.Errorf(err.Error())
	}

	a, err := NewAudio(f, 1000, 3, 10)
	if err != nil {
		t.Errorf(err.Error())
	}

	got := a.NumWindow()
	if got != 1670 {
		t.Errorf("NumWindow() = %d; want 132", got)
	}
}
