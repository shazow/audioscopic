package main

import "github.com/mjibson/go-dsp/spectral"

type Visualizer interface {
	Push(samples []float32)
	SetFreq(float64)
}

type basic struct {
	freq float64
}

func (vis *basic) Push(samples []float32) {
	po := &spectral.PwelchOptions{
		NFFT: 64,
	}
	powers, _ := spectral.Pwelch(float32To64(samples), vis.freq, po)
	_ = powers
	//vis.Set(float64ToInt(powers, 10))
}

func (vis *basic) SetFreq(freq float64) {
	vis.freq = freq
}

func BasicVisualizer() Visualizer {
	return &basic{
		freq: 64,
	}
}
