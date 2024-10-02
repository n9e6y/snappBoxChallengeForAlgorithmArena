package models

type FareEstimate struct {
	DeliveryID int64   `csv:"id_delivery"`
	Fare       float64 `csv:"fare_estimate"`
}
