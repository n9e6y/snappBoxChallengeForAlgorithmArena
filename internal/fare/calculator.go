package fare

import (
	"SBCFAA/internal/models"
	"SBCFAA/internal/processing"
)

func CalculateFare(segments []processing.Segment) models.FareEstimate {
	if len(segments) == 0 {
		return models.FareEstimate{Fare: MinimumFare}
	}

	totalFare := FlagAmount
	deliveryID := segments[0].Start.ID

	for _, segment := range segments {
		if segment.Speed > MovingThreshold {
			totalFare += calculateMovingFare(segment)
		} else {
			totalFare += calculateIdleFare(segment)
		}
	}

	if totalFare < MinimumFare {
		totalFare = MinimumFare
	}

	return models.FareEstimate{
		DeliveryID: deliveryID,
		Fare:       totalFare,
	}
}

func calculateMovingFare(segment processing.Segment) float64 {
	for _, rule := range Rules {
		if rule.State == "MOVING" && rule.Condition(segment.Start.Timestamp) {
			return rule.Amount * segment.Distance
		}
	}
	return 0 // Should never reach here if rules are properly defined
}

func calculateIdleFare(segment processing.Segment) float64 {
	for _, rule := range Rules {
		if rule.State == "IDLE" && rule.Condition(segment.Start.Timestamp) {
			return rule.Amount * segment.Duration.Hours()
		}
	}
	return 0 // Should never reach here if rules are properly defined
}
