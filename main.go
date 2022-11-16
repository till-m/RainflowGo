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
	start := flag.Int("s", 20, "Starting point for damage calculation.")

	fmt.Println("Running...")
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
		stressTemp, err = strconv.ParseFloat(scanner.Text(), 64)
		if err != nil {
			panic(err)
		}
		stress = append(stress, stressTemp)
	}

	// Print to console and write to the outfile
	out, err := os.Create(*outpath)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := out.Close()
		if err != nil {
			panic(err)
		}
	}()

	// Create new writer to write results to file
	w := bufio.NewWriter(out)

	for i := 0; i < *start; i++ {
		outputResults(0, w)
	}

	for i := *start + 1; i <= len(stress); i++ {
		var damage = CalculateDamage(stress[:i], *stressrange)

		outputResults(damage, w)
	}

	w.Flush()
	fmt.Println("\tFinished...")
}

func outputResults(damage float64, w *bufio.Writer) {
	_, err := fmt.Fprintf(w, "%e\n", damage)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%e\n", damage)
}
