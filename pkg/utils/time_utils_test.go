package utils

import (
	"testing"
	"time"
)

func TestIsNightTime(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected bool
	}{
		{
			name:     "Midnight",
			time:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "Just before 5 AM",
			time:     time.Date(2023, 1, 1, 4, 59, 59, 0, time.UTC),
			expected: true,
		},
		{
			name:     "5 AM",
			time:     time.Date(2023, 1, 1, 5, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "Noon",
			time:     time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "11:59 PM",
			time:     time.Date(2023, 1, 1, 23, 59, 59, 0, time.UTC),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNightTime(tt.time)
			if result != tt.expected {
				t.Errorf("IsNightTime(%v) = %v; want %v", tt.time, result, tt.expected)
			}
		})
	}
}

func TestParseTimestamp(t *testing.T) {
	tests := []struct {
		name      string
		timestamp int64
		expected  time.Time
	}{
		{
			name:      "Unix epoch",
			timestamp: 0,
			expected:  time.Unix(0, 0).UTC(),
		},
		{
			name:      "Positive timestamp",
			timestamp: 1609459200, // 2021-01-01 00:00:00
			expected:  time.Unix(1609459200, 0).UTC(),
		},
		{
			name:      "Negative timestamp",
			timestamp: -1000000, // 1969-12-20 13:13:20
			expected:  time.Unix(-1000000, 0).UTC(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseTimestamp(tt.timestamp)
			// Convert result to UTC for comparison
			resultUTC := result.UTC()
			if !resultUTC.Equal(tt.expected) {
				t.Errorf("ParseTimestamp(%d) = %v; want %v", tt.timestamp, resultUTC, tt.expected)
			}
		})
	}
}

func TestCalculateDuration(t *testing.T) {
	tests := []struct {
		name     string
		start    int64
		end      int64
		expected time.Duration
	}{
		{
			name:     "Same time",
			start:    1000,
			end:      1000,
			expected: 0,
		},
		{
			name:     "One hour difference",
			start:    1000,
			end:      4600,
			expected: time.Hour,
		},
		{
			name:     "Negative duration",
			start:    4600,
			end:      1000,
			expected: -time.Hour,
		},
		{
			name:     "30 minutes difference",
			start:    1000,
			end:      2800,
			expected: 30 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateDuration(tt.start, tt.end)
			if result != tt.expected {
				t.Errorf("CalculateDuration(%d, %d) = %v; want %v", tt.start, tt.end, result, tt.expected)
			}
		})
	}
}
