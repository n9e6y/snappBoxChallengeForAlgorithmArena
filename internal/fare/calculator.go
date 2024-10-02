package fare

import (
	"SBCFAA/internal/ingestion"
	"SBCFAA/pkg/utils"
	"math"
	"time"
)

const (
	FlagCharge           = 1.30
	MinimumFare          = 3.47
	MovingRateDay        = 0.74  // per km
	MovingRateNight      = 1.30  // per km
	IdleRate             = 11.90 // per hour
	MovingSpeedThreshold = 10.0  // km/h
	NightStartHour       = 0
	NightEndHour         = 5
)

type FareEstimate struct {
	DeliveryID int64
	Fare       float64
}

func CalculateFares(points []ingestion.DeliveryPoint) []FareEstimate {
	fareEstimates := make(map[int64]float64)

	for i := 1; i < len(points); i++ {
		prev := points[i-1]
		curr := points[i]

		if prev.ID != curr.ID {
			// New delivery, apply flag charge
			fareEstimates[curr.ID] = FlagCharge
			continue
		}

		distance := utils.HaversineDistance(prev.Lat, prev.Lng, curr.Lat, curr.Lng)
		duration := time.Duration(curr.Timestamp-prev.Timestamp) * time.Second
		speed := calculateSpeed(distance, duration)

		fare := calculateSegmentFare(distance, duration, speed, time.Unix(curr.Timestamp, 0))
		fareEstimates[curr.ID] += fare
	}

	return consolidateFares(fareEstimates)
}

func calculateSpeed(distance float64, duration time.Duration) float64 {
	hours := duration.Hours()
	if hours == 0 {
		return 0
	}
	return distance / hours // km/h
}

func calculateSegmentFare(distance float64, duration time.Duration, speed float64, timestamp time.Time) float64 {
	if speed <= MovingSpeedThreshold {
		// Idle state
		return IdleRate * duration.Hours()
	}

	// Moving state
	rate := MovingRateDay
	if isNightTime(timestamp) {
		rate = MovingRateNight
	}
	return rate * distance
}

func isNightTime(t time.Time) bool {
	hour := t.Hour()
	return hour >= NightStartHour && hour < NightEndHour
}

func consolidateFares(fareEstimates map[int64]float64) []FareEstimate {
	var results []FareEstimate
	for id, fare := range fareEstimates {
		// Apply minimum fare
		if fare < MinimumFare {
			fare = MinimumFare
		}
		results = append(results, FareEstimate{
			DeliveryID: id,
			Fare:       math.Round(fare*100) / 100, // Round to 2 decimal places
		})
	}
	return results
}
