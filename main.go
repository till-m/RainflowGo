package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {

	// Flag handling
	inpath := flag.String("i", "", "Path to input file containing line seperated data")
	outpath := flag.String("o", "", "Path to output file where results will be stored")
	stressrange := flag.Float64("r", 10.0, "The range of values that will be counted in the fatigue count")

	flag.Parse()

	// Open the file and create the scanner
	file, err := os.Open(*inpath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// Append to a slice
	var stress []float64
	var stressTemp float64

	for scanner.Scan() {
		stressTemp, _ = strconv.ParseFloat(scanner.Text(), 64)
		stress = append(stress, stressTemp)
	}

	// Remove non-peaks from raw stress data
	stripped := Peaks(stress)
	var n_peaks int = stripped.Len()
	// Perform Rainflow count to get half anf full counts
	half, full := RainflowCounting(stripped)

	// Get the counts of each
	result := GetCounts(half, full, *stressrange)

	// Write results to console
	fmt.Println("Rainflow counter ASTM E1049 85 cl 5.4.4")
	fmt.Println("----------------------------------------")
	fmt.Printf("Input file:\t\t%v\n", *inpath)
	fmt.Printf("Bin size:\t\t%.3f\n", *stressrange)
	fmt.Printf("Data points:\t\t%v\n", len(stress))
	fmt.Printf("Peaks and troughs:\t%v\n\n", n_peaks)

	fmt.Println("--------------------------------------------------------------------------")

	fmt.Printf("Bin Low\t\tBin High\tBin Mean\tRange Mean\tCount\n")

	// Print to console and write to the outfile
	out, err := os.Create(*outpath)
	defer out.Close()

	// Create new writer to write results to file
	w := bufio.NewWriter(out)
	fmt.Fprintf(w, "Bin Low,Bin High,Bin Mean,Range Mean,Count\n")

	var count, meanRange, meanBin float64
	for _, k := range result {
		meanRange, count = k.RangeMeanCount()
		meanBin = k.BinMean()

		if count > 0 {
			fmt.Printf("%.2f\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.2f\n", k.Low, k.High, meanBin, meanRange, count)
			fmt.Fprintf(w, "%.5f,%.5f,%.5f,%.5f,%.5f\n", k.Low, k.High, meanBin, meanRange, count)
		}
	}
	w.Flush()

	fmt.Println("--------------------------------------------------------------------------")
	fmt.Printf("Results written to file %v\n", *outpath)
}
