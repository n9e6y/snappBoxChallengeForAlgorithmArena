package fare

import (
	"math"
	"sync"
	"time"

	"SBCFAA/internal/models"
	"SBCFAA/pkg/utils"
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
	workerPoolSize       = 5
)

func CalculateFares(deliveries <-chan []models.DeliveryPoint) <-chan models.FareEstimate {
	estimatesChan := make(chan models.FareEstimate, 100)

	go func() {
		defer close(estimatesChan)

		var wg sync.WaitGroup
		for i := 0; i < workerPoolSize; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for delivery := range deliveries {
					estimate := calculateFareForDelivery(delivery)
					estimatesChan <- estimate
				}
			}()
		}

		wg.Wait()
	}()

	return estimatesChan
}

func calculateFareForDelivery(delivery []models.DeliveryPoint) models.FareEstimate {
	if len(delivery) == 0 {
		return models.FareEstimate{}
	}

	totalFare := FlagCharge
	for i := 1; i < len(delivery); i++ {
		prevPoint := delivery[i-1]
		currentPoint := delivery[i]

		distance := utils.HaversineDistance(prevPoint.Latitude, prevPoint.Longitude, currentPoint.Latitude, currentPoint.Longitude)
		duration := currentPoint.Timestamp.Sub(prevPoint.Timestamp)
		speed := utils.CalculateSpeed(prevPoint, currentPoint)

		fare := calculateSegmentFare(distance, duration, speed, currentPoint.Timestamp)
		totalFare += fare
	}

	if totalFare < MinimumFare {
		totalFare = MinimumFare
	}

	return models.FareEstimate{
		DeliveryID: delivery[0].ID,
		Fare:       math.Round(totalFare*100) / 100, // Round to 2 decimal places
	}
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

func consolidateFares(estimates map[int64]float64) []models.FareEstimate {
	var results []models.FareEstimate
	for id, fare := range estimates {
		if fare < MinimumFare {
			fare = MinimumFare
		}
		results = append(results, models.FareEstimate{
			DeliveryID: id,
			Fare:       math.Round(fare*100) / 100, // Round to 2 decimal places
		})
	}
	return results
}
