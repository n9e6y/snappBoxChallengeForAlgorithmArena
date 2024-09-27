package models

type FareEstimate struct {
	DeliveryID string  `csv:"id_delivery"`
	Fare       float64 `csv:"fare_estimate"`
}
