package models

import "time"

type DeliveryPoint struct {
	ID        int64     `csv:"id_delivery"`
	Latitude  float64   `csv:"lat"`
	Longitude float64   `csv:"lng"`
	Timestamp time.Time `csv:"timestamp"`
}
