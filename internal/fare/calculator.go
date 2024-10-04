package fare

import (
	"SBCFAA/pkg/utils"
	"math"
	"sync"
	"time"

	"SBCFAA/internal/models"
)

const (
	FlagCharge           = 1.30
	MinimumFare          = 3.47
	MovingRateDay        = 0.74  // per km
	MovingRateNight      = 1.30  // per km
	IdleRate             = 11.90 // per hour
	MovingSpeedThreshold = 10.0  // km/hour
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
		Fare:       math.Round(totalFare*100) / 100, // Round to 2decimal
	}
}

func calculateSegmentFare(distance float64, duration time.Duration, speed float64, timestamp time.Time) float64 {
	if speed <= MovingSpeedThreshold { // Idle state
		return IdleRate * duration.Hours()
	}

	rate := MovingRateDay // Moving state
	if isNightTime(timestamp) {
		rate = MovingRateNight
	}
	return rate * distance
}

func isNightTime(t time.Time) bool {
	hour := t.Hour()
	return hour >= NightStartHour && hour < NightEndHour
}
