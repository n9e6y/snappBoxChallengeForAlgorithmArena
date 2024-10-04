package fare

import (
	"SBCFAA/internal/models"
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestCalculateFares(t *testing.T) {
	createPoint := func(id int64, lat, lon float64, timestamp time.Time) models.DeliveryPoint {
		return models.DeliveryPoint{
			ID:        id,
			Latitude:  lat,
			Longitude: lon,
			Timestamp: timestamp,
		}
	}

	// Test cases
	testCases := []struct {
		name     string
		input    [][]models.DeliveryPoint
		expected []models.FareEstimate
	}{
		{
			name: "Single delivery - day rate",
			input: [][]models.DeliveryPoint{
				{
					createPoint(1, 40.7128, -74.0060, time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)),
					createPoint(1, 40.7128, -74.0070, time.Date(2023, 1, 1, 12, 10, 0, 0, time.UTC)),
				},
			},
			expected: []models.FareEstimate{
				{DeliveryID: 1, Fare: 3.47}, // Minimum fare applied
			},
		},
		{
			name: "Single delivery - night rate",
			input: [][]models.DeliveryPoint{
				{
					createPoint(2, 40.7128, -74.0060, time.Date(2023, 1, 1, 2, 0, 0, 0, time.UTC)),
					createPoint(2, 40.7128, -74.0070, time.Date(2023, 1, 1, 2, 10, 0, 0, time.UTC)),
				},
			},
			expected: []models.FareEstimate{
				{DeliveryID: 2, Fare: 3.47}, // Minimum fare applied
			},
		},
		{
			name: "Single delivery - idle state",
			input: [][]models.DeliveryPoint{
				{
					createPoint(3, 40.7128, -74.0060, time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)),
					createPoint(3, 40.7128, -74.0061, time.Date(2023, 1, 1, 12, 30, 0, 0, time.UTC)),
				},
			},
			expected: []models.FareEstimate{
				{DeliveryID: 3, Fare: 7.25}, // 1.30 (flag) + 5.95 (30 min idle) = 7.25
			},
		},
		{
			name: "Multiple deliveries - mixed scenarios",
			input: [][]models.DeliveryPoint{
				{
					createPoint(4, 40.7128, -74.0060, time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)),
					createPoint(4, 40.7138, -74.0070, time.Date(2023, 1, 1, 12, 10, 0, 0, time.UTC)),
				},
				{
					createPoint(5, 40.7128, -74.0060, time.Date(2023, 1, 1, 2, 0, 0, 0, time.UTC)),
					createPoint(5, 40.7138, -74.0070, time.Date(2023, 1, 1, 2, 10, 0, 0, time.UTC)),
				},
				{
					createPoint(6, 40.7128, -74.0060, time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)),
					createPoint(6, 40.7128, -74.0061, time.Date(2023, 1, 1, 13, 0, 0, 0, time.UTC)),
				},
			},
			expected: []models.FareEstimate{
				{DeliveryID: 4, Fare: 3.47},
				{DeliveryID: 5, Fare: 3.47},
				{DeliveryID: 6, Fare: 13.20}, // 1.30 (flag) + 11.90 (1 hour idle) = 13.20
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create input channel
			deliveriesChan := make(chan []models.DeliveryPoint, len(tc.input))
			for _, delivery := range tc.input {
				deliveriesChan <- delivery
			}
			close(deliveriesChan)

			resultChan := CalculateFares(deliveriesChan) // Call CalculateFares
			// Collect results
			var results []models.FareEstimate
			for estimate := range resultChan {
				results = append(results, estimate)
			}

			// Sort results and expected for consistent comparison
			sort.Slice(results, func(i, j int) bool {
				return results[i].DeliveryID < results[j].DeliveryID
			})
			sort.Slice(tc.expected, func(i, j int) bool {
				return tc.expected[i].DeliveryID < tc.expected[j].DeliveryID
			})

			if !reflect.DeepEqual(results, tc.expected) { // Compare results
				t.Errorf("CalculateFares() = %v, want %v", results, tc.expected)
			}
		})
	}
}

func TestCalculateSegmentFare(t *testing.T) {
	tests := []struct {
		name      string
		distance  float64
		duration  time.Duration
		speed     float64
		timestamp time.Time
		expected  float64
	}{
		{
			name:      "Idle state during day",
			distance:  0.5,
			duration:  30 * time.Minute,
			speed:     5.0,
			timestamp: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			expected:  IdleRate * 0.5, // 30 minutes = 0.5 hours
		},
		{
			name:      "Moving state during day",
			distance:  10.0,
			duration:  30 * time.Minute,
			speed:     20.0,
			timestamp: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			expected:  MovingRateDay * 10.0,
		},
		{
			name:      "Moving state during night",
			distance:  10.0,
			duration:  30 * time.Minute,
			speed:     20.0,
			timestamp: time.Date(2023, 1, 1, 2, 0, 0, 0, time.UTC),
			expected:  MovingRateNight * 10.0,
		},
		{
			name:      "Exactly at speed threshold",
			distance:  5.0,
			duration:  30 * time.Minute,
			speed:     MovingSpeedThreshold,
			timestamp: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			expected:  IdleRate * 0.5, // 30 minutes = 0.5 hours
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateSegmentFare(tt.distance, tt.duration, tt.speed, tt.timestamp)
			if result != tt.expected {
				t.Errorf("calculateSegmentFare() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsNightTime(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected bool
	}{
		{"Day time", time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC), false},
		{"Night time start", time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), true},
		{"Night time end", time.Date(2023, 1, 1, 4, 59, 59, 0, time.UTC), true},
		{"Day time start", time.Date(2023, 1, 1, 5, 0, 0, 0, time.UTC), false},
		{"Just before night time", time.Date(2023, 1, 1, 23, 59, 59, 0, time.UTC), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNightTime(tt.time)
			if result != tt.expected {
				t.Errorf("isNightTime() = %v, want %v", result, tt.expected)
			}
		})
	}
}
