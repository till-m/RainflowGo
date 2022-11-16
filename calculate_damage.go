package main

import "math"

func CalculateDamage(stress []float64, stressrange float64) float64 {

	// Remove non-peaks from raw stress data
	stripped := Peaks(stress)

	// Perform Rainflow count to get half anf full counts
	half, full := RainflowCounting(stripped)

	// Get the counts of each
	result := GetCounts(half, full, stressrange)

	// Slope of the curve
	const m = 8

	// number of cycle for knee point
	const Nk = 2e6

	// Ultimate stress
	const Rm = 865

	// [MPa] Endurance Limit
	const sigaf = Rm / 2

	var damage float64 = 0
	var count, meanRange, meanBin float64
	for _, k := range result {
		meanRange, count = k.RangeMeanCount()
		meanBin = k.BinMean()

		var siga = (meanRange * Rm) / (Rm - meanBin) // Mean stress correction (Goodman)
		damage += math.Pow((count/Nk)*(siga/sigaf), m)
	}

	return damage
}
