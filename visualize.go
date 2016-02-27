package main

import "github.com/mjibson/go-dsp/spectral"

func Visualizer(rate int, channels int) *visualizer {
	return &visualizer{
		channels: channels,
		rate:     rate,
	}
}

type visualizer struct {
	channels int
	rate     int
	prev     []float64
}

func (vis *visualizer) Push(samples []float32) {
	po := &spectral.PwelchOptions{
		Scale_off: true,
	}

	samplesPerChan := int(len(samples) / vis.channels)
	channels := make([][]float64, 0, vis.channels)

	for i := 0; i < vis.channels; i++ {
		channels = append(channels, make([]float64, 0, samplesPerChan))
	}
	for i, s := range samples {
		c := i % vis.channels
		channels[c] = append(channels[c], float64(s))
	}

	// TODO: Handle both channels
	powers, _ := spectral.Pwelch(channels[0], float64(vis.rate), po)

	/*
		for i, p := range powers {
			powers[i] = math.Pow(p, -1)
		}*/

	//fmt.Println(vis.channels, ranges)

	vis.prev = powers[:(len(powers) / 2)]
}

func (vis *visualizer) Sample() []float64 {
	// TODO: Adjust for human frequency sensitivity (http://www.lafavre.us/sound-loudness.jpg)
	return vis.prev
}
