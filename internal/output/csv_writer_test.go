package output

import (
	"SBCFAA/internal/models"
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestWriteCSV(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "csv_writer_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file path
	testFile := filepath.Join(tempDir, "test_output.csv")

	// Create test data
	testData := []models.FareEstimate{
		{DeliveryID: 1, Fare: 10.50},
		{DeliveryID: 2, Fare: 15.75},
		{DeliveryID: 3, Fare: 8.25},
	}

	// Create a channel and send test data
	estimatesChan := make(chan models.FareEstimate, len(testData))
	for _, estimate := range testData {
		estimatesChan <- estimate
	}
	close(estimatesChan)

	// Call the function under test
	err = WriteCSV(testFile, estimatesChan)
	if err != nil {
		t.Fatalf("WriteCSV failed: %v", err)
	}

	// Verify the output file
	file, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("Failed to open output file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	expectedLines := []string{
		"id_delivery,fare_estimate",
		"1,10.50",
		"2,15.75",
		"3,8.25",
	}

	for scanner.Scan() {
		if lineCount >= len(expectedLines) {
			t.Fatalf("Too many lines in output file")
		}
		if scanner.Text() != expectedLines[lineCount] {
			t.Errorf("Line %d: expected '%s', got '%s'", lineCount, expectedLines[lineCount], scanner.Text())
		}
		lineCount++
	}

	if lineCount != len(expectedLines) {
		t.Errorf("Expected %d lines, got %d", len(expectedLines), lineCount)
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("Error reading output file: %v", err)
	}
}

func TestWriteCSVEmptyInput(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "csv_writer_test_empty")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "test_output_empty.csv")

	estimatesChan := make(chan models.FareEstimate)
	close(estimatesChan)

	err = WriteCSV(testFile, estimatesChan)
	if err != nil {
		t.Fatalf("WriteCSV failed with empty input: %v", err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expectedContent := "id_delivery,fare_estimate\n"
	if string(content) != expectedContent {
		t.Errorf("Expected file content to be '%s', got '%s'", expectedContent, string(content))
	}
}

func TestWriteCSVLargeDataset(t *testing.T) {
	// Define the size of the large dataset
	const datasetSize = 1000000 // 1 million records

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "large_test*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Create a channel to send fare estimates
	estimatesChan := make(chan models.FareEstimate, bufferSize)

	// Start a goroutine to generate and send fare estimates
	go func() {
		defer close(estimatesChan)
		for i := 0; i < datasetSize; i++ {
			estimatesChan <- models.FareEstimate{
				DeliveryID: int64(i + 1),
				Fare:       float64(i%1000) + 0.99, // Vary the fare to avoid repetition
			}
		}
	}()

	// Measure the time taken to write the CSV
	startTime := time.Now()

	// Call WriteCSV
	err = WriteCSV(tmpfile.Name(), estimatesChan)
	if err != nil {
		t.Fatalf("WriteCSV failed: %v", err)
	}

	duration := time.Since(startTime)
	t.Logf("Time taken to write %d records: %v", datasetSize, duration)

	// Verify the file contents
	file, err := os.Open(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to open the written file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
		if lineCount == 1 {
			// Check header
			if scanner.Text() != "id_delivery,fare_estimate" {
				t.Errorf("Incorrect header: got %s, want id_delivery,fare_estimate", scanner.Text())
			}
		} else {
			// Check a sample of lines (e.g., every 100,000th line)
			if (lineCount-1)%100000 == 0 {
				parts := strings.Split(scanner.Text(), ",")
				if len(parts) != 2 {
					t.Errorf("Incorrect line format at line %d: %s", lineCount, scanner.Text())
					continue
				}
				id, err := strconv.ParseInt(parts[0], 10, 64)
				if err != nil {
					t.Errorf("Failed to parse ID at line %d: %v", lineCount, err)
				}
				fare, err := strconv.ParseFloat(parts[1], 64)
				if err != nil {
					t.Errorf("Failed to parse fare at line %d: %v", lineCount, err)
				}
				expectedFare := float64((id-1)%1000) + 0.99
				if fare != expectedFare {
					t.Errorf("Incorrect fare at line %d: got %f, want %f", lineCount, fare, expectedFare)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	expectedLineCount := datasetSize + 1 // +1 for header
	if lineCount != expectedLineCount {
		t.Errorf("Incorrect number of lines: got %d, want %d", lineCount, expectedLineCount)
	}

	// Calculate and log the write speed
	fileStat, err := file.Stat()
	if err != nil {
		t.Fatalf("Failed to get file stats: %v", err)
	}
	fileSize := fileStat.Size()
	speedMBPerSecond := float64(fileSize) / duration.Seconds() / 1024 / 1024

	t.Logf("File size: %d bytes", fileSize)
	t.Logf("Write speed: %.2f MB/s", speedMBPerSecond)
}
