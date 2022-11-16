package main

import (
	"container/list"
	"math"
	"sort"
)

// Peaks takes a list of the raw stress values and removes all
// intermediate values which are not peaks or troughs
func Peaks(s []float64) *list.List {

	// Output slice
	var stripped = list.New()

	// Value of the change in stress
	var ds float64

	// Loop over the stress ranges
	for i, v := range s {

		// If its the first value, append it
		if i == 0 {
			stripped.PushBack(v)
		} else {

			// if ds is 0 cant divide by
			if ds == 0 {
				ds = v - s[i-1]
				continue
			}

			// If the current ds is a different sign
			if (v-s[i-1])/ds < 0 {
				stripped.PushBack(s[i-1])
				ds = v - s[i-1]
			}
		}
	}
	// Append the final value
	stripped.PushBack(s[len(s)-1])

	return stripped
}

// RainflowCounting takes a set of peaks only - should be processed by peaks()
// returns a slice of half cycles and whole cycles
// as per ASTM E1049 85 Cl 5.4.4
func RainflowCounting(p *list.List) ([]float64, []float64) {
	var e_i = p.Front()
	var X, Y float64

	// Slices of half and full to append ranges to
	var half, full []float64

	// b can be changed to true to break out of the loop if the conditions are met
	var b bool = false

	for {
		if b == false {
			// (1) - Read the next peak of valley
			var e_i1 = e_i.Next()
			var e_i2 = e_i1.Next()
			if e_i2 == nil {
				b = true
				continue
			}

			Y = e_i1.Value.(float64) - e_i.Value.(float64)
			X = e_i2.Value.(float64) - e_i1.Value.(float64)

			// (3a) X < Y
			if math.Abs(X) < math.Abs(Y) {
				// go to (1)
				e_i = e_i.Next()
				continue
			} else {
				// X >= Y go to (4)
				// If range Y contains the starting point S
				if p.Front() == e_i {

					// go to (5) - Count Y as a half cycle drop S
					half = append(half, math.Abs(Y))
					p.Remove(p.Front())

					e_i = p.Front()
				} else {
					// (4) Remove the peak and the valley of Y and count as full cycle
					full = append(full, math.Abs(Y))
					p.Remove(e_i.Next())
					p.Remove(e_i)

					e_i = p.Front()
				}
			}
		} else {
			// Collect remaining cycles and attribute them to half cycles
			for true {
				e_i = p.Front()
				if p.Len() == 1 {
					break
				}
				half = append(half, math.Abs(e_i.Next().Value.(float64)-e_i.Value.(float64)))
				p.Remove(e_i)
			}

			return half, full
		}
		// Break conditions
		if p.Len() < 3 {
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
