package main

import (
	"math"
	"sort"
)

// Peaks takes a list of the raw stress values and removes all
// intermediate values which are not peaks or troughs
func Peaks(s []float64) []float64 {

	// Output slice
	var stripped []float64

	// Value of the change in stress
	var ds float64

	// Loop over the stress ranges
	for i, v := range s {

		// If its the first value, append it
		if i == 0 {
			stripped = append(stripped, v)
		} else {

			// if ds is 0 cant divide by
			if ds == 0 {
				ds = v - s[i-1]
				continue
			}

			// If the current ds is a different sign
			if (v-s[i-1])/ds < 0 {
				stripped = append(stripped, s[i-1])
				ds = v - s[i-1]
			}
		}
	}
	// Append the final value
	stripped = append(stripped, s[len(s)-1])

	return stripped
}

// RainflowCounting takes a set of peaks only - should be processed by peaks()
// returns a slice of half cycles and whole cycles
// as per ASTM E1049 85 Cl 5.4.4
func RainflowCounting(p []float64) ([]float64, []float64) {

	var X, Y float64
	var i int = 0

	// Slices of half and full to append ranges to
	var half, full []float64

	// b can be changed to true to break out of the loop if the conditions are met
	var b bool = false

	for {
		if b == false {
			// (1) - Read the next peak of valley
			Y = p[i+1] - p[i]
			X = p[i+2] - p[i+1]

			// (3a) X < Y
			if math.Abs(X) < math.Abs(Y) {
				// go to (1)
				i++
			} else {
				// X >= Y go to (4)
				// If range Y contains the starting point S
				if i == 0 {

					// go to (5) - Count Y as a half cycle drop S
					half = append(half, math.Abs(Y))
					p = removeElement(p, 0)

					i = 0
				} else {
					// (4) Remove the peak and the valley of Y and count as full cycle
					full = append(full, math.Abs(Y))
					p = removeElement(p, i)
					p = removeElement(p, i)
					i = 0
				}
			}
		} else {
			// Collect reamaining cycles and attribute them to half cycles
			for i := range p {
				if i == len(p)-1 {
					break
				}

				half = append(half, math.Abs(p[i+1]-p[i]))
			}
			return half, full
		}
		// Break conditions
		if len(p) < 3 {
			b = true
		}
		if i > len(p)-3 {
			b = true
		}
	}

}

// GetCounts takes the slice of half count ranges, full count ranges and range interval r
// and returns a list of Count structs representing the count of each bin. 
func GetCounts(half []float64, full []float64, r float64) []Count {

	// Sort the half and full slices into ascending order
	sort.Float64s(half)
	sort.Float64s(full)

	// Get the min and max - there could be empty half or full arrays
	var min, max float64
	if (len(half) > 0) && (len(full) > 0) {
		min = math.Min(half[0], full[0])
		max = math.Max(half[len(half)-1], full[len(full)-1])
	} else {
		if len(half) > 0 {
			min = half[0]
			max = half[len(half)-1]
		} else {
			min = full[0]
			max = full[len(full)-1]
		}
	}

	// Get lowest values within the bin
	binLow := math.Floor(min/r) * r
	binHigh := math.Floor(max/r) * r

	// Create the Count objects
	var countSlice []Count
	
	// Set initial bin low and high values
	bL := binLow
	bH := binLow + r
	var c Count
	
	// Generate empty array of count objects
	for {
		c = Count{Low: bL,
			High: bH,
			Half: make([]float64, 0),
			Full: make([]float64, 0)}

		countSlice = append(countSlice, c)

		bL += r
		bH += r

		if bL > binHigh {
			break
		}
	}
	
	// binCounter is incremented as required
	var binCounter int
	
	// Loop over half range
	for _, v := range half {

		// If the value is larger than the bin, keep adding the range until it fits
		for {
			if (v < countSlice[binCounter].High) && (v >= countSlice[binCounter].Low) {
				break
			} else {
				// Move on to next bin
				binCounter++
			}
		}
		// Add half ranges to this bin
		countSlice[binCounter].Half = append(countSlice[binCounter].Half, v)
	}

	// Reset the bin values
	binCounter = 0

	// Loop over ful range
	for _, v := range full {

		// If the value is larger than the bin, keep adding the range until it fits
		for {
			if (v < countSlice[binCounter].High) && (v >= countSlice[binCounter].Low) {
				break
			} else {
				// Move on to next bin
				binCounter++
			}
		}
		// Add full ranges to this bin
		countSlice[binCounter].Full = append(countSlice[binCounter].Full, v)
	}

	return countSlice
}

// removeElements removes an element at index i from a slice of type []float64
func removeElement(s []float64, i int) []float64 {
	// Remove the element at index i from a.
	copy(s[i:], s[i+1:]) // Shift a[i+1:] left one index.
	s[len(s)-1] = 0      // Erase last element (write zero value).
	s = s[:len(s)-1]     // Truncate slice.
	return s
}

// GetMeanCount takes the list of counts created by GetCounts and returns
// a map of range means and total counts
func GetMeanCount(c []Count) map[float64]float64 {

	// Initialise bin mean, count and the resulting map
	var binMean, count float64
	result := make(map[float64]float64)

	for _, v := range c {

		// Assert that there are non zero counts
		if (len(v.Half) > 0) || (len(v.Full) > 0) {
			binMean, count = v.RangeMeanCount()

			result[binMean] = count
		}
	}

	return result
}
