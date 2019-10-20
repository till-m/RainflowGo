package main

// Count is a struct that represents the count of fatigue cycles
// Low: Lower bound of ranges
// High: Higher bound of ranges
// Half: Slice of half cycle ranges
// Full: Slice of full cycle ranges
type Count struct {
	Low  float64
	High float64
	Half []float64
	Full []float64
}

// BinMean returns the midpoint between the Low and High bins of Count struct
func (c *Count) BinMean() float64 {
	return (c.Low + c.High) / 2
}

// RangeMeanCount returns the mean weighted average of the ranges of the Count struct and the total count of that bin
func (c *Count) RangeMeanCount() (float64, float64) {
	var count, sumTotal float64

	// Loop over half ranges
	for _, v := range c.Half {
		count += 0.5
		sumTotal += v * 0.5
	}

	// Loop over full ranges
	for _, v := range c.Full {
		count += 1.0
		sumTotal += v
	}

	return sumTotal / count, count
}

// CheckBins checks of all the values in the Half and Full Slices are between the bin Low and High bounds
// Returns true if OK
func (c *Count) CheckBins() bool {

	// Loop over half ranges
	for _, v := range c.Half {
		if (v < c.Low) || (v > c.High) {
			return false
		}
	}
	// Loop over half ranges
	for _, v := range c.Full {
		if (v < c.Low) || (v > c.High) {
			return false
		}
	}

	return true

}
