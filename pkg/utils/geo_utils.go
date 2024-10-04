package utils

import (
	"SBCFAA/internal/models"
	"math"
)

const (
	earthRadiusKm = 6371.0 // Earth's radius in KM
	maxLatitude   = 90.0
	minLatitude   = -90.0
)

// HaversineDistance calculates the great-circle distance between two points on a sphere
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	if !isValidLatitude(lat1) || !isValidLatitude(lat2) { // Check for invalid latitudes
		return math.NaN()
	}

	// Convert latitude and longitude to radians
	lat1 = toRadians(lat1)
	lon1 = toRadians(lon1)
	lat2 = toRadians(lat2)
	lon2 = toRadians(lon2)

	// Haversine formula
	dLat := lat2 - lat1
	dLon := lon2 - lon1
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

// CalculateSpeed calculates the speed between two delivery points in km/h
func CalculateSpeed(p1, p2 models.DeliveryPoint) float64 {
	distance := HaversineDistance(p1.Latitude, p1.Longitude, p2.Latitude, p2.Longitude)
	duration := p2.Timestamp.Sub(p1.Timestamp).Hours()
	if duration == 0 {
		return math.Inf(1)
	}
	return distance / duration
}

func isValidLatitude(lat float64) bool {
	return lat >= minLatitude && lat <= maxLatitude
}

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
