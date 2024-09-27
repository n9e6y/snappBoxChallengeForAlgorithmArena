package ingestion

import (
	"encoding/csv"
	"io"
	"strconv"
	"time"

	"SBCFAA/internal/models"
)

// ReadCSV reads delivery points from a CSV file
func ReadCSV(reader io.Reader) ([]models.DeliveryPoint, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = 4 // id_delivery, lat, lng, timestamp

	// Read and discard header
	if _, err := csvReader.Read(); err != nil {
		return nil, err
	}

	var points []models.DeliveryPoint
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		point, err := parseDeliveryPoint(record)
		if err != nil {
			return nil, err
		}

		points = append(points, point)
	}

	return points, nil
}

// parseDeliveryPoint converts a slice of strings to a DeliveryPoint
func parseDeliveryPoint(record []string) (models.DeliveryPoint, error) {
	id, err := parseID(record[0])
	if err != nil {
		return models.DeliveryPoint{}, err
	}

	lat, err := parseFloat(record[1])
	if err != nil {
		return models.DeliveryPoint{}, err
	}

	lng, err := parseFloat(record[2])
	if err != nil {
		return models.DeliveryPoint{}, err
	}

	timestamp, err := parseTimestamp(record[3])
	if err != nil {
		return models.DeliveryPoint{}, err
	}

	return models.DeliveryPoint{
		ID:        id,
		Latitude:  lat,
		Longitude: lng,
		Timestamp: timestamp,
	}, nil
}

// parseID converts a string to an int64
func parseID(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// parseFloat converts a string to a float64
func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// parseTimestamp converts a string (assumed to be Unix timestamp) to time.Time
func parseTimestamp(s string) (time.Time, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(i, 0), nil
}
