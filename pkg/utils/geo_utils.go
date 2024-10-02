package utils

import (
	"SBCFAA/internal/models"
	"math"
)

const (
	earthRadiusKm = 6371.0 // Earth's radius in kilometers
)

// HaversineDistance calculates the great-circle distance between two points on a sphere
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	lat1 = toRadians(lat1)
	lon1 = toRadians(lon1)
	lat2 = toRadians(lat2)
	lon2 = toRadians(lon2)

	dlat := lat2 - lat1
	dlon := lon2 - lon1
	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// CalculateSpeed calculates the speed between two delivery points in km/h
func CalculateSpeed(p1, p2 models.DeliveryPoint) float64 {
	distance := HaversineDistance(p1.Latitude, p1.Longitude, p2.Latitude, p2.Longitude)
	duration := p2.Timestamp.Sub(p1.Timestamp).Hours()
	if duration == 0 {
		return 0
	}
	return distance / duration
}
