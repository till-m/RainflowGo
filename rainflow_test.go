package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestRemoveElements(t *testing.T) {
	testSlice := []float64{0, 1, 2, 3, 4, 5}
	resultSlice := removeElement(testSlice, 2)

	expected := []float64{0, 1, 3, 4, 5}

	if reflect.DeepEqual(resultSlice, expected) == false {
		t.Errorf("removeElement failed to remove correct element")
	}
}

func TestPeaks(t *testing.T) {
	testSlice := []float64{0, 1, 2, 3, 3, 2, 1, 0}
	resultSlice := Peaks(testSlice)

	expected := []float64{0, 3, 0}

	if reflect.DeepEqual(resultSlice, expected) == false {
		t.Errorf("TestPeaks failed to find peaks")
	}
}

func TestRainflowCounter(t *testing.T) {
	testSlice := []float64{-2, 1, -3, 5, -1, 3, -4, 4, -2}
	resultHalf, resultFull := RainflowCounting(testSlice)

	expectedHalf := []float64{3, 4, 6, 8, 8, 9}
	expectedFull := []float64{4}

	sort.Float64s(resultHalf)
	sort.Float64s(resultFull)
	sort.Float64s(expectedHalf)
	sort.Float64s(expectedFull)

	if (reflect.DeepEqual(resultHalf, expectedHalf) != true) || (reflect.DeepEqual(resultFull, expectedFull) != true) {
		t.Errorf("Rainflow counter failed test to solve example from Figure 6 of ASTM E1049 85\n")
		t.Errorf("Should be...%v\n%v\n", expectedHalf, expectedFull)
		t.Errorf("Got \n%v\n%v\n", resultHalf, resultFull)
	}
}

func TestGetMeanCounts(t *testing.T) {
	testSlice := []float64{-2, 1, -3, 5, -1, 3, -4, 4, -2}
	resultHalf, resultFull := RainflowCounting(testSlice)
	counts := GetCounts(resultHalf, resultFull, 1.0)
	result := GetMeanCount(counts)

	expected := map[float64]float64{
		3: 0.5,
		4: 1.5,
		6: 0.5,
		8: 1,
		9: 0.5,
	}

	if reflect.DeepEqual(result, expected) == false {
		t.Errorf("GetCounts failed test to solve example from Figure 6 of ASTM E1049 85\n")
		t.Errorf("Should be %v\n", expected)
		t.Errorf("Got %v\n", result)
	}
}
