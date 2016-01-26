package main

import (
	"math"
	"math/cmplx"
)

func float32To64(a []float32) []float64 {
	r := make([]float64, len(a))
	for n, v := range a {
		r[n] = float64(v)
	}
	return r
}

func float64ToInt(a []float64, mult float64) []int {
	r := make([]int, len(a))
	for n, v := range a {
		r[n] = 100 + int(math.Log10(v)*mult)
	}
	return r
}

func f32ToComplex(a []float32) []complex128 {
	r := make([]complex128, len(a))
	for n, v := range a {
		r[n] = complex(float64(v), 0)
	}
	return r
}

func complexToInt(a []complex128) []int {
	r := make([]int, len(a))
	for n, v := range a {
		// Convert complex phase to percentage int
		r[n] = int(((cmplx.Phase(v) + math.Pi) / math.Pi) * 100)
	}
	return r
}
