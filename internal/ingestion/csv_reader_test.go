package ingestion

import (
	"SBCFAA/internal/models"
	"io/ioutil"
	"math"
	"os"
	"testing"
	"time"
)

func TestReadAndFilterCSV(t *testing.T) {
	// Test cases
	testCases := []struct {
		name               string
		input              string
		expectedDeliveries int
		expectedErrors     int
	}{
		{
			name: "Valid input",
			input: `id,lat,lng,timestamp
1,40.7128,-74.0060,1609459200
1,40.7129,-74.0061,1609459260
2,40.7130,-74.0062,1609459320
2,40.7131,-74.0063,1609459380`,
			expectedDeliveries: 2, // Two separate deliveries
			expectedErrors:     0,
		},
		{
			name: "Invalid speed filtered out",
			input: `id,lat,lng,timestamp
1,40.7128,-74.0060,1609459200
1,50.7129,-84.0061,1609459260`, // Large distance in short time
			expectedDeliveries: 1, // One delivery with one point
			expectedErrors:     0,
		},
		{
			name: "Invalid data format",
			input: `id,lat,lng,timestamp
1,invalid,invalid,invalid`,
			expectedDeliveries: 0,
			expectedErrors:     1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary file
			tmpfile, err := ioutil.TempFile("", "test.csv")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpfile.Name())

			// Write test data to the file
			if _, err := tmpfile.Write([]byte(tc.input)); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatalf("Failed to close temp file: %v", err)
			}

			// Call the function
			pointsChan, errChan := ReadAndFilterCSV(tmpfile.Name())

			// Count the deliveries and errors
			deliveryCount := 0
			for range pointsChan {
				deliveryCount++
			}

			errorCount := 0
			for range errChan {
				errorCount++
			}

			if deliveryCount != tc.expectedDeliveries { // Check results
				t.Errorf("Expected %d deliveries, got %d", tc.expectedDeliveries, deliveryCount)
			}
			if errorCount != tc.expectedErrors {
				t.Errorf("Expected %d errors, got %d", tc.expectedErrors, errorCount)
			}
		})
	}
}

func TestReadAndFilterCSVEndToEnd(t *testing.T) {
	input := `id,lat,lng,timestamp
1,40.7128,-74.0060,1609459200
1,40.7129,-74.0061,1609459260
1,40.7130,-74.0062,1609459320
2,40.7131,-74.0063,1609459380
2,40.7132,-74.0064,1609459440
2,50.7133,-84.0065,1609459500
3,40.7134,-74.0066,1609459560
3,40.7135,-74.0067,1609459620`

	// Create a temporary file
	tmpfile, err := ioutil.TempFile("", "test_e2e.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	// Write test data to the file
	if _, err := tmpfile.Write([]byte(input)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Call the function
	pointsChan, errChan := ReadAndFilterCSV(tmpfile.Name())

	// Process the results
	deliveries := make(map[int64][]models.DeliveryPoint)
	for points := range pointsChan {
		if len(points) > 0 {
			deliveries[points[0].ID] = points
		}
	}

	// Check for errors
	for err := range errChan {
		t.Errorf("Unexpected error: %v", err)
	}

	// Validate the results
	expectedDeliveries := 3
	if len(deliveries) != expectedDeliveries {
		t.Errorf("Expected %d deliveries, got %d", expectedDeliveries, len(deliveries))
	}

	// Check specific deliveries
	if len(deliveries[1]) != 3 {
		t.Errorf("Expected 3 points for delivery 1, got %d", len(deliveries[1]))
	}
	if len(deliveries[2]) != 2 {
		t.Errorf("Expected 2 points for delivery 2, got %d", len(deliveries[2]))
	}
	if len(deliveries[3]) != 2 {
		t.Errorf("Expected 2 points for delivery 3, got %d", len(deliveries[3]))
	}

	// Validate filtering (the point with high speed should be filtered out)
	for _, points := range deliveries[2] {
		if points.Latitude > 50 {
			t.Errorf("High-speed point was not filtered out: %v", points)
		}
	}
}

func TestParseDeliveryPoint(t *testing.T) {
	tests := []struct {
		name          string
		input         []string
		expected      models.DeliveryPoint
		expectedError bool
	}{
		{
			name:  "Valid input",
			input: []string{"1", "40.7128", "-74.0060", "1609459200"},
			expected: models.DeliveryPoint{
				ID:        1,
				Latitude:  40.7128,
				Longitude: -74.0060,
				Timestamp: time.Unix(1609459200, 0),
			},
			expectedError: false,
		},
		{
			name:          "Invalid ID",
			input:         []string{"invalid", "40.7128", "-74.0060", "1609459200"},
			expectedError: true,
		},
		{
			name:          "Invalid latitude",
			input:         []string{"1", "invalid", "-74.0060", "1609459200"},
			expectedError: true,
		},
		{
			name:          "Invalid longitude",
			input:         []string{"1", "40.7128", "invalid", "1609459200"},
			expectedError: true,
		},
		{
			name:          "Invalid timestamp",
			input:         []string{"1", "40.7128", "-74.0060", "invalid"},
			expectedError: true,
		},
		{
			name:          "Too few fields",
			input:         []string{"1", "40.7128", "-74.0060"},
			expectedError: true,
		},
		{
			name:          "Too many fields",
			input:         []string{"1", "40.7128", "-74.0060", "1609459200", "extra"},
			expectedError: true,
		},
		{
			name:  "Edge case: zero values",
			input: []string{"0", "0", "0", "0"},
			expected: models.DeliveryPoint{
				ID:        0,
				Latitude:  0,
				Longitude: 0,
				Timestamp: time.Unix(0, 0),
			},
			expectedError: false,
		},
		{
			name:  "Edge case: negative ID",
			input: []string{"-1", "40.7128", "-74.0060", "1609459200"},
			expected: models.DeliveryPoint{
				ID:        -1,
				Latitude:  40.7128,
				Longitude: -74.0060,
				Timestamp: time.Unix(1609459200, 0),
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseDeliveryPoint(tt.input)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if result.ID != tt.expected.ID {
					t.Errorf("Expected ID %d, but got %d", tt.expected.ID, result.ID)
				}

				if result.Latitude != tt.expected.Latitude {
					t.Errorf("Expected Latitude %f, but got %f", tt.expected.Latitude, result.Latitude)
				}

				if result.Longitude != tt.expected.Longitude {
					t.Errorf("Expected Longitude %f, but got %f", tt.expected.Longitude, result.Longitude)
				}

				if !result.Timestamp.Equal(tt.expected.Timestamp) {
					t.Errorf("Expected Timestamp %v, but got %v", tt.expected.Timestamp, result.Timestamp)
				}
			}
		})
	}
}

func TestCalculateSpeed(t *testing.T) {
	tests := []struct {
		name     string
		p1       models.DeliveryPoint
		p2       models.DeliveryPoint
		expected float64
	}{
		{
			name: "Same location, 1 hour apart",
			p1: models.DeliveryPoint{
				Latitude:  40.7128,
				Longitude: -74.0060,
				Timestamp: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			p2: models.DeliveryPoint{
				Latitude:  40.7128,
				Longitude: -74.0060,
				Timestamp: time.Date(2023, 1, 1, 13, 0, 0, 0, time.UTC),
			},
			expected: 0,
		},
		{
			name: "Different location, 1 hour apart",
			p1: models.DeliveryPoint{
				Latitude:  40.7128,
				Longitude: -74.0060,
				Timestamp: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			p2: models.DeliveryPoint{
				Latitude:  40.7138,
				Longitude: -74.0070,
				Timestamp: time.Date(2023, 1, 1, 13, 0, 0, 0, time.UTC),
			},
			expected: 0.1396, // Approximate speed in km/h
		},
		{
			name: "Same timestamp",
			p1: models.DeliveryPoint{
				Latitude:  40.7128,
				Longitude: -74.0060,
				Timestamp: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			p2: models.DeliveryPoint{
				Latitude:  40.7138,
				Longitude: -74.0070,
				Timestamp: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateSpeed(tt.p1, tt.p2)
			if tt.expected == 0 {
				if result != 0 {
					t.Errorf("calculateSpeed() = %v, want %v", result, tt.expected)
				}
			} else {
				tolerance := 0.0001 // Adjust this value based on your required precision
				if math.Abs(result-tt.expected) > tolerance {
					t.Errorf("calculateSpeed() = %v, want %v (within %v)", result, tt.expected, tolerance)
				}
			}
		})
	}
}
