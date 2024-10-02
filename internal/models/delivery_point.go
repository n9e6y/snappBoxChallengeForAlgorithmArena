package models

import "time"

type DeliveryPoint struct {
	ID        int64
	Latitude  float64
	Longitude float64
	Timestamp time.Time
}

func NewDeliveryPoint(id int64, lat, lng float64, timestamp int64) DeliveryPoint {
	return DeliveryPoint{
		ID:        id,
		Latitude:  lat,
		Longitude: lng,
		Timestamp: time.Unix(timestamp, 0),
	}
}
