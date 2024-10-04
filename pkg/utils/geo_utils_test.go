package utils

import (
	"SBCFAA/internal/models"
	"math"
	"testing"
	"time"
)

func TestCalculateSpeed(t *testing.T) {
	tests := []struct {
		name     string
		p1       models.DeliveryPoint
		p2       models.DeliveryPoint
		expected float64
	}{
		{
			name: "Same point",
			p1: models.DeliveryPoint{
				Latitude:  40.7128,
				Longitude: -74.0060,
				Timestamp: time.Unix(1609459200, 0),
			},
			p2: models.DeliveryPoint{
				Latitude:  40.7128,
				Longitude: -74.0060,
				Timestamp: time.Unix(1609459260, 0),
			},
			expected: 0,
		},
		{
			name: "Different points",
			p1: models.DeliveryPoint{
				Latitude:  40.7128,
				Longitude: -74.0060,
				Timestamp: time.Unix(1609459200, 0),
			},
			p2: models.DeliveryPoint{
				Latitude:  40.7129,
				Longitude: -74.0061,
				Timestamp: time.Unix(1609459260, 0),
			},
			expected: 0.8371704438758566, // This is the correct value based on the Haversine formula
		},
		{
			name: "Zero time difference",
			p1: models.DeliveryPoint{
				Latitude:  40.7128,
				Longitude: -74.0060,
				Timestamp: time.Unix(1609459200, 0),
			},
			p2: models.DeliveryPoint{
				Latitude:  40.7129,
				Longitude: -74.0061,
				Timestamp: time.Unix(1609459200, 0),
			},
			expected: math.Inf(1), // Infinite speed due to zero time difference
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateSpeed(tt.p1, tt.p2)
			if math.Abs(result-tt.expected) > 0.0000001 { // Using a very small delta for float comparison
				t.Errorf("CalculateSpeed() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHaversineDistance(t *testing.T) {
	tests := []struct {
		name      string
		lat1      float64
		lon1      float64
		lat2      float64
		lon2      float64
		expected  float64
		tolerance float64
	}{
		{
			name:      "Same point",
			lat1:      40.7128,
			lon1:      -74.0060,
			lat2:      40.7128,
			lon2:      -74.0060,
			expected:  0,
			tolerance: 0.1,
		},
		{
			name:      "New York to Los Angeles",
			lat1:      40.7128,
			lon1:      -74.0060,
			lat2:      34.0522,
			lon2:      -118.2437,
			expected:  3935.746,
			tolerance: 1.0,
		},
		{
			name:      "London to Tokyo",
			lat1:      51.5074,
			lon1:      -0.1278,
			lat2:      35.6762,
			lon2:      139.6503,
			expected:  9559.784,
			tolerance: 5.0,
		},
		{
			name:      "North Pole to South Pole",
			lat1:      90,
			lon1:      0,
			lat2:      -90,
			lon2:      0,
			expected:  20015.087,
			tolerance: 10.0,
		},
		{
			name:      "Antipodes",
			lat1:      51.5074,
			lon1:      -0.1278,
			lat2:      -51.5074,
			lon2:      179.8722,
			expected:  20015.087,
			tolerance: 10.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HaversineDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("HaversineDistance() = %v, want %v (tolerance: %v)", result, tt.expected, tt.tolerance)
			}
		})
	}
}

func TestHaversineDistanceEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		lat1     float64
		lon1     float64
		lat2     float64
		lon2     float64
		expected float64
	}{
		{
			name:     "Invalid latitude (>90)",
			lat1:     91,
			lon1:     0,
			lat2:     0,
			lon2:     0,
			expected: math.NaN(),
		},
		{
			name:     "Invalid latitude (<-90)",
			lat1:     0,
			lon1:     0,
			lat2:     -91,
			lon2:     0,
			expected: math.NaN(),
		},
		{
			name:     "Longitude wrap-around (positive)",
			lat1:     0,
			lon1:     179,
			lat2:     0,
			lon2:     -179,
			expected: 222.4, // km (approximately)
		},
		{
			name:     "Longitude wrap-around (negative)",
			lat1:     0,
			lon1:     -179,
			lat2:     0,
			lon2:     179,
			expected: 222.4, // km (approximately)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HaversineDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			if math.IsNaN(tt.expected) {
				if !math.IsNaN(result) {
					t.Errorf("HaversineDistance() = %v, want NaN", result)
				}
			} else if math.Abs(result-tt.expected) > 0.1 { // Using a larger tolerance for edge cases
				t.Errorf("HaversineDistance() = %v, want %v", result, tt.expected)
			}
		})
	}
}
