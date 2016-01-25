package main

import "math"

type Visualizer interface {
	Push(samples []float32)
	Render()
}

type basic struct {
	value float64
	Set   func(float64)
}

func (vis *basic) Push(samples []float32) {
	// root mean square based on https://github.com/mdlayher/waveform/blob/master/samplereducefunc.go#L18
	// TODO: FFT and stuff.
	var sumSquare float64
	for _, s := range samples {
		sumSquare += math.Pow(float64(s), 2)
	}
	vis.value = math.Sqrt(sumSquare / float64(len(samples)))
	vis.Set(vis.value)
}

func (vis *basic) Render() {
}

func BasicVisualizer(setter func(float64)) Visualizer {
	return &basic{
		Set: setter,
	}
}
