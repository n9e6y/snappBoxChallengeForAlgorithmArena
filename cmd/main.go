package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

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

	startTime := time.Now()

	// Read and filter input data
	log.Println("Reading and filtering input data...")
	pointsChan, errChan := ingestion.ReadAndFilterCSV(*inputFile)

	// Calculate fares
	log.Println("Calculating fares...")
	estimatesChan := fare.CalculateFares(pointsChan)

	// Write results to CSV
	log.Println("Writing results to CSV...")
	if err := output.WriteCSV(*outputFile, estimatesChan); err != nil {
		log.Fatalf("Error writing output data: %v", err)
	}

	// Check for any errors from reading/filtering
	for err := range errChan {
		log.Printf("Error during processing: %v", err)
	}

	duration := time.Since(startTime)
	log.Printf("Fare estimation completed successfully in %v. Results written to %s\n", duration, *outputFile)

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
