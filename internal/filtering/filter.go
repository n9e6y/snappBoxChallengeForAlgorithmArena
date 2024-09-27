package filtering

import (
	"SBCFAA/internal/models"
	"SBCFAA/pkg/utils"
)

const maxSpeed = 100.0 //Km

// FilterInvalidPoints removes points that result in unrealistic speeds
func FilterInvalidPoints(points []models.DeliveryPoint) []models.DeliveryPoint {
	if len(points) < 2 {
		return points
	}

	filtered := []models.DeliveryPoint{points[0]}

	for i := 1; i < len(points); i++ {
		p1 := filtered[len(filtered)-1]
		p2 := points[i]

		if isValidSegment(p1, p2) {
			filtered = append(filtered, p2)
		}
	}

	return filtered
}

func isValidSegment(p1, p2 models.DeliveryPoint) bool {
	distance := utils.HaversineDistance(p1.Latitude, p1.Longitude, p2.Latitude, p2.Longitude)
	duration := p2.Timestamp.Sub(p1.Timestamp).Hours()

	if duration == 0 {
		return false // Avoid division by zero
	}

	speed := distance / duration

	return speed <= maxSpeed
}
