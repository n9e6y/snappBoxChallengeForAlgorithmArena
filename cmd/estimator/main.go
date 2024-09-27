package main

import (
	"SBCFAA/internal/fare"
	"SBCFAA/internal/filtering"
	"SBCFAA/internal/ingestion"
	"SBCFAA/internal/models"
	"SBCFAA/internal/output"
	"SBCFAA/internal/processing"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: estimator <input_file> <output_file>")
	}

	inputFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to open input file: %v", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outputFile.Close()

	// Step 1: Read input data
	points, err := ingestion.ReadCSV(inputFile)
	if err != nil {
		log.Fatalf("Failed to read input data: %v", err)
	}

	// Step 2: Filter invalid points
	filteredPoints := filtering.FilterInvalidPoints(points)

	// Step 3: Process segments
	segments := processing.ProcessSegments(filteredPoints)

	// Step 4: Calculate fares
	estimates := make([]models.FareEstimate, 0, len(segments))
	for _, segment := range segments {
		estimate := fare.CalculateFare([]processing.Segment{segment})
		estimates = append(estimates, estimate)
	}

	// Step 5: Write output
	if err := output.WriteCSV(outputFile, estimates); err != nil {
		log.Fatalf("Failed to write output: %v", err)
	}

	log.Println("Fare estimation completed successfully")
}
