package main

import "github.com/mjibson/go-dsp/spectral"

func Analyze(samples []float32, rate int) []float64 {
	po := &spectral.PwelchOptions{
		NFFT:      16,
		Scale_off: true,
	}

	powers, ranges := spectral.Pwelch(float32To64(samples), float64(rate), po)

	sum := 0.0
	for i, p := range powers {
		sum += p
		if p <= 0.001 {
			powers[i] = 0.001
		}
	}
	_ = ranges
	//fmt.Printf("[%f~%f] %f = %v\n", ranges[0], ranges[len(ranges)-1], sum, powers[:8])

	return powers[:3]
}
