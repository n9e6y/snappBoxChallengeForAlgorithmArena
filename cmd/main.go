package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"SBCFAA/internal/fare"
	"SBCFAA/internal/ingestion"
	"SBCFAA/internal/output"
)

func main() {
	// Define command-line flags
	inputFile := flag.String("input", "", "Input CSV file path")
	outputFile := flag.String("output", "fare_estimates.csv", "Output CSV file path")
	cpuProfile := flag.String("cpuprofile", "", "Write cpu profile to file")
	memProfile := flag.String("memprofile", "", "Write memory profile to file")
	flag.Parse()

	// Check if input file is provided
	if *inputFile == "" {
		log.Fatal("Please provide an input file using the -input flag")
	}

	// CPU profiling
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal("Could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("Could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	// Read and filter input data
	fmt.Println("Reading and filtering input data...")
	filteredData, err := ingestion.ReadAndFilterCSV(*inputFile)
	if err != nil {
		log.Fatalf("Error reading input data: %v", err)
	}

	// Calculate fares
	fmt.Println("Calculating fares...")
	fareEstimates := fare.CalculateFares(filteredData)

	// Write results to CSV
	fmt.Println("Writing results to CSV...")
	err = output.WriteCSV(*outputFile, fareEstimates)
	if err != nil {
		log.Fatalf("Error writing output data: %v", err)
	}

	fmt.Printf("Fare estimation completed successfully. Results written to %s\n", *outputFile)

	// Memory profiling
	if *memProfile != "" {
		f, err := os.Create(*memProfile)
		if err != nil {
			log.Fatal("Could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC() // Get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("Could not write memory profile: ", err)
		}
	}
}
