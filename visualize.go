package main

import "github.com/mjibson/go-dsp/spectral"

type visualizer struct {
	rate int
	prev []float64
}

func (vis *visualizer) Push(samples []float32) {
	po := &spectral.PwelchOptions{
		NFFT:      16,
		Scale_off: true,
	}

	powers, ranges := spectral.Pwelch(float32To64(samples), float64(vis.rate), po)

	sum := 0.0
	for i, p := range powers {
		sum += p
		if p <= 0.001 {
			powers[i] = 0.001
		}
	}
	_ = ranges
	//fmt.Printf("[%f~%f] %f = %v\n", ranges[0], ranges[len(ranges)-1], sum, powers[:8])

	vis.prev = powers
}

func (vis *visualizer) Sample() []float64 {
	// TODO: Adjust for human frequency sensitivity (http://www.lafavre.us/sound-loudness.jpg)
	return vis.prev
}
