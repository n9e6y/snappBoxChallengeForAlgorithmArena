package utils

import (
	"time"
)

// IsNightTime checks if the given time is within night hours (00:00 to 05:00)
func IsNightTime(t time.Time) bool {
	hour := t.Hour()
	return hour >= 0 && hour < 5
}

// ParseTimestamp converts a Unix timestamp to a time.Time
func ParseTimestamp(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

// CalculateDuration returns the duration between two timestamps
func CalculateDuration(start, end int64) time.Duration {
	return time.Duration(end-start) * time.Second
}
