package utils

import (
	"math"
)

const (
	earthRadiusKm = 6371.0
)

// HaversineDistance calculates the great-circle distance between two points on a sphere
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := toRadians(lat2 - lat1)
	dLon := toRadians(lon2 - lon1)
	lat1 = toRadians(lat1)
	lat2 = toRadians(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}

func toRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}
